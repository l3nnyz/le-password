@echo off
setlocal
cd /d "%~dp0"

rem -trimpath : retire les chemins locaux du binaire
rem -s -w     : retire la table des symboles et les infos de debug (DWARF)
set CGO_ENABLED=0

echo Compilation de le-password.exe...
go build -trimpath -ldflags="-s -w" -o le-password.exe .
if errorlevel 1 (
    echo Echec de la compilation.
    exit /b 1
)

for %%F in (le-password.exe) do echo OK : le-password.exe ^(%%~zF octets^)
endlocal
