#!/bin/bash
# Script para auto-corregir problemas comunes de Markdown
# Usa markdownlint-cli para aplicar correcciones automáticas

set -e

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔧 Auto-fixing Markdown files...${NC}"

# Verificar si npx está disponible
if ! command -v npx &> /dev/null; then
    echo -e "${RED}❌ Error: npx no está instalado${NC}"
    echo -e "${BLUE}💡 Instala Node.js/npm primero${NC}"
    exit 1
fi

# Ejecutar markdownlint-cli con auto-fix
echo -e "${BLUE}📝 Ejecutando markdownlint-cli --fix...${NC}"

# Patrones a incluir
PATTERNS=(
    "*.md"
    "docs/**/*.md"
    ".github/**/*.md"
    ".vscode/**/*.md"
)

# Patrones a excluir
IGNORE_PATTERNS=(
    "node_modules"
    "vendor"
    "volumes"
    "tmp"
    "reference-repos"
    ".next"
    "dist"
    "build"
)

# Construir comando con exclusiones
IGNORE_ARGS=""
for pattern in "${IGNORE_PATTERNS[@]}"; do
    IGNORE_ARGS="$IGNORE_ARGS --ignore $pattern"
done

# Ejecutar markdownlint con --fix
echo -e "${YELLOW}Corrigiendo archivos Markdown...${NC}"

# Usar npx para ejecutar markdownlint-cli sin instalación global
if npx markdownlint-cli --fix $IGNORE_ARGS '**/*.md' 2>&1; then
    echo ""
    echo -e "${GREEN}✅ Archivos Markdown corregidos exitosamente${NC}"
    echo -e "${GREEN}📝 Los cambios se han aplicado automáticamente${NC}"
    exit 0
else
    EXIT_CODE=$?
    
    # Si el código de salida es 1, significa que hubo problemas pero se corrigieron algunos
    if [ $EXIT_CODE -eq 1 ]; then
        echo ""
        echo -e "${YELLOW}⚠️  Algunos archivos tenían problemas pero se intentaron corregir${NC}"
        echo -e "${BLUE}💡 Revisa los cambios con 'git diff'${NC}"
        exit 0
    else
        echo ""
        echo -e "${RED}❌ Error ejecutando markdownlint${NC}"
        exit $EXIT_CODE
    fi
fi
