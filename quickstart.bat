@echo off
REM YouTube 视频分析工作流 - Windows 快速启动脚本

echo ========================================
echo GoWorkFlow YouTube 视频分析 - 快速启动
echo ========================================
echo.

REM 步骤 1: 检查 Go
echo [1/5] 检查系统依赖...
echo ----------------------------
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [X] Go 未安装
    echo 请访问 https://golang.org/dl/ 安装 Go 1.16+
    pause
    exit /b 1
)
echo [√] Go 已安装

REM 检查 MySQL
where mysql >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [X] MySQL 未安装
    echo 请安装 MySQL 8.0+ 或使用 Docker
    pause
    exit /b 1
)
echo [√] MySQL 已安装

echo.

REM 步骤 2: 检查配置文件
echo [2/5] 检查配置文件...
echo ----------------------------
if not exist .env (
    echo 创建 .env 文件...
    copy config\.env.example .env >nul
    echo [√] .env 文件已创建
) else (
    echo [√] .env 文件已存在
)

echo.

REM 读取配置
set DB_HOST=localhost
set DB_PORT=3306
set DB_USER=root
set DB_NAME=workflow

set /p DB_PASSWORD="请输入 MySQL root 密码: "

echo.

REM 步骤 3: 初始化数据库
echo [3/5] 初始化数据库...
echo ----------------------------

REM 创建数据库
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% -e "CREATE DATABASE IF NOT EXISTS %DB_NAME% CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [X] 无法连接数据库，请检查密码和权限
    pause
    exit /b 1
)
echo [√] 数据库已创建

REM 运行迁移
echo 运行数据库迁移...
for %%f in (migrations\*.sql) do (
    echo   ^> 执行 %%~nxf
    mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% %DB_NAME% < "%%f" 2>nul
)
echo [√] 数据库迁移完成

echo.

REM 步骤 4: 构建应用
echo [4/5] 构建应用...
echo ----------------------------

if not exist bin\workflow-api.exe (
    echo 首次构建应用...
    go build -o bin\workflow-api.exe cmd\workflow-api\main.go
    if %ERRORLEVEL% NEQ 0 (
        echo [X] 构建失败
        pause
        exit /b 1
    )
) else (
    set /p REBUILD="应用已存在，是否重新构建? [y/N]: "
    if /i "%REBUILD%"=="y" (
        go build -o bin\workflow-api.exe cmd\workflow-api\main.go
    )
)
echo [√] 应用构建完成

echo.

REM 步骤 5: 创建必要目录
echo [5/5] 创建必要目录...
echo ----------------------------

if not exist reports mkdir reports
if not exist web\css mkdir web\css
if not exist web\js mkdir web\js
echo [√] 目录已创建

echo.
echo ========================================
echo [√] 所有准备工作已完成！
echo ========================================
echo.

REM 更新 .env 文件中的密码
powershell -Command "(gc .env) -replace 'DB_PASSWORD=.*', 'DB_PASSWORD=%DB_PASSWORD%' | Out-File -encoding ASCII .env"

echo 选择启动方式:
echo   1) 立即启动（前台运行）
echo   2) 后台运行（新窗口）
echo   3) 仅显示启动命令
echo.

set /p CHOICE="请选择 [1-3]: "

if "%CHOICE%"=="1" (
    echo.
    echo [启动应用] 按 Ctrl+C 停止...
    echo.
    bin\workflow-api.exe
) else if "%CHOICE%"=="2" (
    echo.
    echo [后台启动]...
    start "GoWorkFlow" bin\workflow-api.exe
    echo [√] 应用已在新窗口启动
    echo.
    echo 访问地址: http://localhost:8080
    echo.
) else if "%CHOICE%"=="3" (
    echo.
    echo 手动启动命令:
    echo.
    echo   bin\workflow-api.exe
    echo.
) else (
    echo 无效选择
    pause
    exit /b 1
)

echo.
echo 快速使用指南:
echo ----------------------------
echo 1. 访问 Web 界面: http://localhost:8080
echo 2. 点击 '流程管理' -^> 'YouTube 视频智能分析'
echo 3. 输入 YouTube 视频地址
echo 4. 等待 1-3 分钟，查看作业进度
echo 5. 查看生成的报告
echo.
echo 详细文档: YOUTUBE_ANALYSIS_GUIDE.md
echo.
echo 祝您使用愉快！
pause
