#!/bin/bash

# Script para crear una biblioteca de medios simulada para desarrollo/testing
# Crea archivos vacíos que simulan películas y series

set -e

MEDIA_DIR="./volumes/media-library/library"
MOVIES_DIR="$MEDIA_DIR/movies"
TVSHOWS_DIR="$MEDIA_DIR/tv"

echo "🎬 Creando biblioteca de medios simulada..."

# Crear directorios si no existen
mkdir -p "$MOVIES_DIR"
mkdir -p "$TVSHOWS_DIR"

# Función para crear un archivo sparse de tamaño específico
create_file() {
    local filepath="$1"
    local size_mb="$2"
    
    # Crear el directorio padre si no existe
    mkdir -p "$(dirname "$filepath")"
    
    # Crear archivo sparse (instantáneo, no ocupa espacio real en disco)
    truncate -s "${size_mb}M" "$filepath"
}

echo ""
echo "📽️  Creando películas..."

# Películas de ejemplo con tamaños realistas
declare -A movies=(
    ["The Matrix (1999)/The Matrix (1999) - 1080p.mkv"]=2048
    ["Inception (2010)/Inception (2010) - 1080p.mkv"]=2560
    ["The Shawshank Redemption (1994)/The Shawshank Redemption (1994) - 1080p.mkv"]=1920
    ["Pulp Fiction (1994)/Pulp Fiction (1994) - 1080p.mkv"]=2304
    ["The Dark Knight (2008)/The Dark Knight (2008) - 1080p.mkv"]=2688
    ["Forrest Gump (1994)/Forrest Gump (1994) - 1080p.mkv"]=2176
    ["Fight Club (1999)/Fight Club (1999) - 1080p.mkv"]=2432
    ["Goodfellas (1990)/Goodfellas (1990) - 1080p.mkv"]=2240
    ["The Godfather (1972)/The Godfather (1972) - 1080p.mkv"]=2816
    ["The Lord of the Rings The Fellowship of the Ring (2001)/The Fellowship of the Ring (2001) - 1080p.mkv"]=3072
)

for movie in "${!movies[@]}"; do
    size=${movies[$movie]}
    filepath="$MOVIES_DIR/$movie"
    echo "  ✓ Creando: $movie ($size MB)"
    create_file "$filepath" $size
done

echo ""
echo "📺 Creando series..."

# Series de ejemplo con temporadas y episodios
declare -a series=(
    "Breaking Bad"
    "Game of Thrones"
    "The Office"
    "Friends"
    "Stranger Things"
)

for show in "${series[@]}"; do
    echo "  📁 Serie: $show"
    
    # Crear 2-3 temporadas por serie
    num_seasons=$((RANDOM % 2 + 2))
    
    for season in $(seq 1 $num_seasons); do
        season_dir="$TVSHOWS_DIR/$show/Season $(printf %02d $season)"
        mkdir -p "$season_dir"
        
        # Crear 8-12 episodios por temporada
        num_episodes=$((RANDOM % 5 + 8))
        
        for episode in $(seq 1 $num_episodes); do
            episode_file="$season_dir/${show} - S$(printf %02d $season)E$(printf %02d $episode) - 1080p.mkv"
            # Episodios de ~800MB - 1.2GB
            size=$((RANDOM % 400 + 800))
            create_file "$episode_file" $size
        done
        
        echo "    ✓ Temporada $season: $num_episodes episodios"
    done
done

echo ""
echo "📊 Estadísticas de la biblioteca:"
echo ""

# Calcular tamaños
movies_count=$(find "$MOVIES_DIR" -type f -name "*.mkv" 2>/dev/null | wc -l)
movies_size=$(du -sh "$MOVIES_DIR" 2>/dev/null | cut -f1)

tvshows_count=$(find "$TVSHOWS_DIR" -type f -name "*.mkv" 2>/dev/null | wc -l)
tvshows_size=$(du -sh "$TVSHOWS_DIR" 2>/dev/null | cut -f1)

total_size=$(du -sh "$MEDIA_DIR/library" 2>/dev/null | cut -f1)

echo "  Películas: $movies_count archivos ($movies_size)"
echo "  Series: $tvshows_count episodios ($tvshows_size)"
echo "  Total: $total_size"

echo ""
echo "✅ Biblioteca de medios simulada creada exitosamente!"
echo ""
echo "📍 Ubicación: $MEDIA_DIR/"
echo ""
echo "💡 Próximos pasos:"
echo "   1. Agregar estas rutas en Radarr y Sonarr:"
echo "      - Radarr: /media-library/library/movies"
echo "      - Sonarr: /media-library/library/tv"
echo "   2. Escanear la biblioteca en cada servicio"
echo "   3. Ejecutar sincronización en KeeperCheky"
echo ""
