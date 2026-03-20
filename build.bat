@echo off
echo Building Razer Battery Monitor...
go build -ldflags="-H windowsgui" -o RazerBatteryMonitor.exe .
if %ERRORLEVEL% == 0 (
    echo Build successful: RazerBatteryMonitor.exe
) else (
    echo Build failed.
    pause
)
