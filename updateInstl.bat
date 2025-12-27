@echo off
setlocal enabledelayedexpansion

for /r %%f in (*.go) do (
    findstr /C:"github.com/yourcompany/thirdparty-sdk/" "%%f" >nul 2>&1
    if !errorlevel! equ 0 (
        powershell -Command "(Get-Content '%%f' -Raw) -replace 'github.com/yourcompany/thirdparty-sdk/', 'github.com/yassine-manai/go_zr_sdk/' | Set-Content '%%f' -NoNewline"
        echo Updated: %%f
    )
)

echo Import update complete!
