@echo off
REM Build and run all services (Windows)
SETLOCAL ENABLEDELAYEDEXPANSION
SET ROOT=%~dp0
IF NOT EXIST "%ROOT%logs" MKDIR "%ROOT%logs"

ECHO Building Go services...

:: Inventory
necho Building inventory...
pushd "%ROOT%services\inventory"
ngo build -o "%ROOT%services\inventory\inventory.exe" . 2> "%ROOT%logs\inventory_build.log"
IF %ERRORLEVEL% NEQ 0 (
    echo Inventory build failed. See %ROOT%logs\inventory_build.log
) ELSE (
    echo Inventory built OK
)
popd

:: Orders
necho Building orders...
pushd "%ROOT%services\orders"
go build -o "%ROOT%services\orders\orders.exe" . 2> "%ROOT%logs\orders_build.log"
IF %ERRORLEVEL% NEQ 0 (
    echo Orders build failed. See %ROOT%logs\orders_build.log
) ELSE (
    echo Orders built OK
)
popd

ECHO Stopping previous instances if running...
:: attempt to stop previous binaries by name
tasklist /FI "IMAGENAME eq inventory.exe" 2>NUL | find /I "inventory.exe" >NUL && (taskkill /IM inventory.exe /F >NUL 2>&1)
tasklist /FI "IMAGENAME eq orders.exe" 2>NUL | find /I "orders.exe" >NUL && (taskkill /IM orders.exe /F >NUL 2>&1)

ECHO Starting services in new windows...

:: Start inventory
start "inventory" "%ROOT%services\inventory\inventory.exe" > "%ROOT%logs\inventory_run.log" 2>&1
ECHO Inventory log: %ROOT%logs\inventory_run.log

:: Start orders
start "orders" "%ROOT%services\orders\orders.exe" > "%ROOT%logs\orders_run.log" 2>&1
ECHO Orders log: %ROOT%logs\orders_run.log

:: Start frontend in new window (install deps then dev)
start "frontend" cmd /K "cd /d "%ROOT%project" && npm install --no-audit --no-fund > "%ROOT%logs\frontend_install.log" 2>&1 && npm run dev > "%ROOT%logs\frontend_run.log" 2>&1"
ECHO Frontend log: %ROOT%logs\frontend_run.log

ECHO All done. Check logs folder for build/run outputs.
ENDLOCAL
pause

