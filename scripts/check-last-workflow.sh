#!/bin/bash

# Script para revisar el estado del último workflow de manera automática
# Requiere: gh CLI instalado y autenticado

REPO="carcheky/keepercheky"

echo "🔍 Revisando último workflow..."
echo ""

# Obtener información del último workflow
LAST_RUN=$(gh run list --repo "$REPO" --limit 1 --json databaseId,status,conclusion,name,headBranch,event,createdAt,url 2>/dev/null)

if [ -z "$LAST_RUN" ] || [ "$LAST_RUN" = "[]" ]; then
    echo "❌ No se encontraron workflows"
    exit 1
fi

# Parsear información
RUN_ID=$(echo "$LAST_RUN" | jq -r '.[0].databaseId')
STATUS=$(echo "$LAST_RUN" | jq -r '.[0].status')
CONCLUSION=$(echo "$LAST_RUN" | jq -r '.[0].conclusion')
NAME=$(echo "$LAST_RUN" | jq -r '.[0].name')
BRANCH=$(echo "$LAST_RUN" | jq -r '.[0].headBranch')
EVENT=$(echo "$LAST_RUN" | jq -r '.[0].event')
CREATED=$(echo "$LAST_RUN" | jq -r '.[0].createdAt')
URL=$(echo "$LAST_RUN" | jq -r '.[0].url')

# Mostrar información básica
echo "📊 Información del Workflow"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Nombre:       $NAME"
echo "  ID:           $RUN_ID"
echo "  Branch:       $BRANCH"
echo "  Evento:       $EVENT"
echo "  Creado:       $CREATED"
echo "  URL:          $URL"
echo ""

# Determinar estado con emojis
if [ "$STATUS" = "completed" ]; then
    if [ "$CONCLUSION" = "success" ]; then
        echo "✅ Estado:      COMPLETADO EXITOSAMENTE"
    elif [ "$CONCLUSION" = "failure" ]; then
        echo "❌ Estado:      FALLÓ"
    elif [ "$CONCLUSION" = "cancelled" ]; then
        echo "🚫 Estado:      CANCELADO"
    elif [ "$CONCLUSION" = "skipped" ]; then
        echo "⏭️  Estado:      OMITIDO"
    else
        echo "❓ Estado:      COMPLETADO ($CONCLUSION)"
    fi
elif [ "$STATUS" = "in_progress" ]; then
    echo "⏳ Estado:      EN PROGRESO"
elif [ "$STATUS" = "queued" ]; then
    echo "⏸️  Estado:      EN COLA"
else
    echo "❓ Estado:      $STATUS"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Si falló, mostrar jobs que fallaron
if [ "$CONCLUSION" = "failure" ]; then
    echo ""
    echo "🔍 Revisando jobs fallidos..."
    echo ""
    
    gh run view "$RUN_ID" --repo "$REPO" --json jobs --jq '.jobs[] | select(.conclusion == "failure") | "  ❌ Job: \(.name)\n     Conclusión: \(.conclusion)\n"' 2>/dev/null
    
    echo ""
    echo "📋 Para ver logs completos:"
    echo "   gh run view $RUN_ID --log"
    echo ""
    echo "🔗 Ver en navegador:"
    echo "   gh run view $RUN_ID --web"
fi

# Si está en progreso, mostrar progreso
if [ "$STATUS" = "in_progress" ]; then
    echo ""
    echo "⏳ Jobs en progreso..."
    echo ""
    
    gh run view "$RUN_ID" --repo "$REPO" --json jobs --jq '.jobs[] | "  \(if .status == "completed" then "✅" elif .status == "in_progress" then "⏳" else "⏸️" end) \(.name) - \(.status)"' 2>/dev/null
    
    echo ""
    echo "🔄 Para monitorear en tiempo real:"
    echo "   gh run watch $RUN_ID"
fi

# Si fue exitoso, mostrar releases/tags creados
if [ "$CONCLUSION" = "success" ] && [ "$NAME" = "Release" ]; then
    echo ""
    echo "🎉 Release exitoso!"
    echo ""
    echo "📦 Últimas releases:"
    gh release list --repo "$REPO" --limit 3 2>/dev/null | while read -r line; do
        echo "   $line"
    done
    
    echo ""
    echo "🏷️  Últimos tags:"
    git tag --sort=-creatordate | head -3 | while read -r tag; do
        echo "   $tag"
    done
    
    echo ""
    echo "🐳 Paquetes Docker:"
    DOCKER_PACKAGES=$(gh api /user/packages/container/keepercheky/versions 2>/dev/null | jq -r 'if type == "array" then .[] | .metadata.container.tags | join(", ") else empty end' 2>/dev/null | head -3)
    if [ -n "$DOCKER_PACKAGES" ]; then
        echo "$DOCKER_PACKAGES" | while read -r tags; do
            echo "   $tags"
        done
    else
        echo "   (No disponible o sin permisos)"
    fi
fi

echo ""
