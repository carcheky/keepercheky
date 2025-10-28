#!/bin/bash

# Script para limpiar completamente todos los releases, tags y packages de GitHub
# Requiere: gh CLI instalado y autenticado

set -e

REPO="carcheky/keepercheky"

echo "🧹 Limpiando todos los releases, tags y packages de GitHub..."
echo ""

# 1. Borrar todos los GitHub Releases
echo "📦 Borrando GitHub Releases..."
gh release list --repo "$REPO" --limit 1000 | while read -r line; do
    tag=$(echo "$line" | awk '{print $1}')
    if [ -n "$tag" ]; then
        echo "  Borrando release: $tag"
        gh release delete "$tag" --repo "$REPO" --yes 2>/dev/null || true
    fi
done
echo "✅ Releases borrados"
echo ""

# 2. Borrar todos los tags remotos
echo "🏷️  Borrando tags remotos..."
git ls-remote --tags origin | awk '{print $2}' | sed 's/refs\/tags\///' | while read -r tag; do
    if [ -n "$tag" ]; then
        echo "  Borrando tag remoto: $tag"
        git push --delete origin "$tag" 2>/dev/null || true
    fi
done
echo "✅ Tags remotos borrados"
echo ""

# 3. Borrar todos los tags locales
echo "🏷️  Borrando tags locales..."
git tag -l | while read -r tag; do
    if [ -n "$tag" ]; then
        echo "  Borrando tag local: $tag"
        git tag -d "$tag" 2>/dev/null || true
    fi
done
echo "✅ Tags locales borrados"
echo ""

# 4. Borrar packages de GitHub Container Registry
echo "📦 Borrando packages de GitHub Container Registry..."
echo "⚠️  NOTA: Los packages deben borrarse manualmente desde:"
echo "   https://github.com/$REPO/pkgs/container/keepercheky"
echo ""
echo "   O con este comando (requiere token con permisos delete:packages):"
echo "   gh api -X DELETE /user/packages/container/keepercheky"
echo ""

# 5. Verificar limpieza
echo "🔍 Verificando limpieza..."
echo ""
echo "Tags locales restantes:"
git tag -l | wc -l | xargs echo "  "
echo ""
echo "Tags remotos restantes:"
git ls-remote --tags origin | wc -l | xargs echo "  "
echo ""
echo "Releases restantes:"
gh release list --repo "$REPO" --limit 10 | wc -l | xargs echo "  "
echo ""

echo "✅ Limpieza completada!"
echo ""
echo "📝 Próximos pasos:"
echo "   1. Verifica que no quedan packages en: https://github.com/$REPO/pkgs/container/keepercheky"
echo "   2. Haz un commit nuevo para que semantic-release cree v1.0.0-dev.1"
