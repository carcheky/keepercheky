#!/usr/bin/env python3
"""
Script para aplicar un parche con formato personalizado.
El formato esperado es:
diff --git a/archivo b/archivo
--- a/archivo
+++ b/archivo
@@ -linea,cantidad +linea,cantidad @@
-texto antiguo
+texto nuevo
"""

import re
import sys
from pathlib import Path

def apply_patch(patch_file, dry_run=False):
    """Aplica un parche personalizado."""
    with open(patch_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Dividir por archivos
    file_sections = re.split(r'^diff --git ', content, flags=re.MULTILINE)[1:]
    
    changes_made = 0
    files_modified = set()
    
    for section in file_sections:
        lines = section.split('\n')
        
        # Extraer nombre de archivo
        first_line = lines[0]
        match = re.match(r'a/(.+?) b/\1', first_line)
        if not match:
            continue
        
        filename = match.group(1)
        filepath = Path(filename)
        
        if not filepath.exists():
            print(f"⚠️  Archivo no encontrado: {filename}")
            continue
        
        # Leer contenido actual
        with open(filepath, 'r', encoding='utf-8') as f:
            file_content = f.read()
        
        original_content = file_content
        
        # Procesar cambios
        i = 0
        while i < len(lines):
            line = lines[i]
            
            # Buscar hunks
            hunk_match = re.match(r'^@@ -(\d+),(\d+) \+(\d+),(\d+) @@', line)
            if hunk_match:
                old_start = int(hunk_match.group(1))
                old_count = int(hunk_match.group(2))
                new_start = int(hunk_match.group(3))
                new_count = int(hunk_match.group(4))
                
                i += 1
                old_text = ""
                new_text = ""
                
                # Recoger líneas del hunk
                while i < len(lines) and not lines[i].startswith('@@'):
                    if lines[i].startswith('-'):
                        old_text += lines[i][1:] + '\n'
                    elif lines[i].startswith('+'):
                        new_text += lines[i][1:] + '\n'
                    i += 1
                
                # Eliminar el último \n si existe
                old_text = old_text.rstrip('\n')
                new_text = new_text.rstrip('\n')
                
                # Aplicar reemplazo
                if old_text in file_content:
                    file_content = file_content.replace(old_text, new_text, 1)
                    changes_made += 1
                else:
                    print(f"⚠️  No se encontró texto en {filename} línea {old_start}")
            else:
                i += 1
        
        # Guardar si hubo cambios
        if file_content != original_content:
            if not dry_run:
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(file_content)
            files_modified.add(filename)
    
    return len(files_modified), changes_made

if __name__ == '__main__':
    patch_file = 'parche.patch'
    dry_run = '--dry-run' in sys.argv
    
    if dry_run:
        print("🔍 Modo dry-run - sin cambios reales")
    
    files, changes = apply_patch(patch_file, dry_run)
    
    if dry_run:
        print(f"\n✅ Se aplicarían cambios en {files} archivos ({changes} modificaciones)")
    else:
        print(f"\n✅ Parche aplicado: {files} archivos modificados ({changes} cambios)")
