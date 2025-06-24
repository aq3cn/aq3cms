@echo off
setlocal enabledelayedexpansion

REM aq3cms Windows 启动脚本
REM 用于在 Windows 系统上快速启动 aq3cms 系统

title aq3cms Startup Script

REM 颜色定义（Windows 10+ 支持 ANSI 颜色）
set "RED=[31m"
set "GREEN=[32m"
set "YELLOW=[33m"
set "BLUE=[34m"
set "NC=[0m"

REM 默认值
set "DEV_MODE=false"
set "USE_DOCKER=false"
set "ACTION="

REM 日志函数
:log_info
echo %GREEN%[INFO]%NC% %~1
goto :eof

:log_warn
echo %YELLOW%[WARN]%NC% %~1
goto :eof

:log_error
echo %RED%[ERROR]%NC% %~1
goto :eof

:log_debug
echo %BLUE%[DEBUG]%NC% %~1
goto :eof

REM 检查命令是否存在
:check_command
where %1 >nul 2>&1
if errorlevel 1 (
    call :log_error "%1 命令未找到，请先安装 %1"
    exit /b 1
)
goto :eof

REM 检查文件是否存在
:check_file
if not exist "%~1" (
    call :log_error "文件 %~1 不存在"
    exit /b 1
)
goto :eof

REM 检查目录是否存在
:check_directory
if not exist "%~1" (
    call :log_warn "目录 %~1 不存在，正在创建..."
    mkdir "%~1"
)
goto :eof

REM 显示帮助信息
:show_help
echo aq3cms Windows 启动脚本
echo.
echo 用法: %~nx0 [选项]
echo.
echo 选项:
echo   -h, --help          显示帮助信息
echo   -d, --dev           启动开发环境
echo   -p, --prod          启动生产环境
echo   -b, --build         构建应用
echo   -t, --test          运行测试
echo   -c, --clean         清理构建文件
echo   --docker            使用 Docker 启动
echo   --docker-dev        使用 Docker 启动开发环境
echo   --stop              停止服务
echo   --restart           重启服务
echo   --logs              查看日志
echo   --health            检查健康状态
echo.
echo 示例:
echo   %~nx0 --dev            # 启动开发环境
echo   %~nx0 --docker         # 使用 Docker 启动生产环境
echo   %~nx0 --docker-dev     # 使用 Docker 启动开发环境
echo   %~nx0 --build          # 构建应用
echo   %~nx0 --test           # 运行测试
goto :eof

REM 检查依赖
:check_dependencies
call :log_info "检查依赖..."

if "%USE_DOCKER%"=="true" (
    call :check_command docker
    call :check_command docker-compose
) else (
    call :check_command go
    where make >nul 2>&1
    if errorlevel 1 (
        call :log_warn "make 命令未找到，将使用 go 命令"
    )
)

call :log_info "依赖检查完成"
goto :eof

REM 初始化环境
:init_environment
call :log_info "初始化环境..."

REM 创建必要的目录
call :check_directory "logs"
call :check_directory "uploads"
call :check_directory "data"

REM 检查配置文件
if not exist "config.yaml" (
    if exist "config.yaml.example" (
        call :log_warn "config.yaml 不存在，正在复制示例配置文件..."
        copy "config.yaml.example" "config.yaml"
    ) else (
        call :log_error "配置文件不存在，请先创建 config.yaml"
        exit /b 1
    )
)

REM 检查环境变量文件
if not exist ".env" (
    if exist ".env.example" (
        call :log_warn ".env 文件不存在，正在复制示例文件..."
        copy ".env.example" ".env"
    )
)

call :log_info "环境初始化完成"
goto :eof

REM 构建应用
:build_app
call :log_info "构建应用..."
where make >nul 2>&1
if errorlevel 1 (
    go build -o bin\aq3cms.exe .\cmd\server
) else (
    make build
)
call :log_info "构建完成"
goto :eof

REM 运行测试
:run_tests
call :log_info "运行测试..."
where make >nul 2>&1
if errorlevel 1 (
    go test -v .\...
) else (
    make test
)
call :log_info "测试完成"
goto :eof

REM 清理构建文件
:clean_build
call :log_info "清理构建文件..."
if exist "bin" rmdir /s /q "bin"
if exist "tmp" rmdir /s /q "tmp"
if exist "coverage.out" del "coverage.out"
if exist "coverage.html" del "coverage.html"
call :log_info "清理完成"
goto :eof

REM 启动开发环境
:start_dev
call :log_info "启动开发环境..."

if "%USE_DOCKER%"=="true" (
    docker-compose -f docker-compose.dev.yml up -d
) else (
    REM 检查 Air 是否安装
    where air >nul 2>&1
    if errorlevel 1 (
        call :log_warn "Air 未安装，正在安装..."
        go install github.com/cosmtrek/air@latest
    )
    
    REM 启动热重载
    air -c .air.toml
)
goto :eof

REM 启动生产环境
:start_prod
call :log_info "启动生产环境..."

if "%USE_DOCKER%"=="true" (
    docker-compose up -d
) else (
    REM 构建应用
    call :build_app
    
    REM 启动应用
    bin\aq3cms.exe
)
goto :eof

REM 停止服务
:stop_services
call :log_info "停止服务..."

if "%USE_DOCKER%"=="true" (
    if "%DEV_MODE%"=="true" (
        docker-compose -f docker-compose.dev.yml down
    ) else (
        docker-compose down
    )
) else (
    REM 停止 aq3cms 进程
    taskkill /f /im aq3cms.exe 2>nul
    taskkill /f /im air.exe 2>nul
)

call :log_info "服务已停止"
goto :eof

REM 重启服务
:restart_services
call :log_info "重启服务..."
call :stop_services
timeout /t 2 /nobreak >nul

if "%DEV_MODE%"=="true" (
    call :start_dev
) else (
    call :start_prod
)
goto :eof

REM 查看日志
:show_logs
call :log_info "查看日志..."

if "%USE_DOCKER%"=="true" (
    if "%DEV_MODE%"=="true" (
        docker-compose -f docker-compose.dev.yml logs -f aq3cms-dev
    ) else (
        docker-compose logs -f aq3cms
    )
) else (
    if exist "logs\app.log" (
        type "logs\app.log"
    ) else (
        call :log_warn "日志文件不存在"
    )
)
goto :eof

REM 健康检查
:health_check
call :log_info "检查健康状态..."

REM 等待服务启动
timeout /t 5 /nobreak >nul

REM 检查健康端点
curl -f http://localhost:8080/health >nul 2>&1
if errorlevel 1 (
    call :log_error "服务健康检查失败"
    exit /b 1
) else (
    call :log_info "服务健康状态正常"
)
goto :eof

REM 解析命令行参数
:parse_args
if "%~1"=="" goto :main_logic
if "%~1"=="-h" goto :help_and_exit
if "%~1"=="--help" goto :help_and_exit
if "%~1"=="-d" (
    set "DEV_MODE=true"
    set "ACTION=start"
    shift
    goto :parse_args
)
if "%~1"=="--dev" (
    set "DEV_MODE=true"
    set "ACTION=start"
    shift
    goto :parse_args
)
if "%~1"=="-p" (
    set "DEV_MODE=false"
    set "ACTION=start"
    shift
    goto :parse_args
)
if "%~1"=="--prod" (
    set "DEV_MODE=false"
    set "ACTION=start"
    shift
    goto :parse_args
)
if "%~1"=="-b" (
    set "ACTION=build"
    shift
    goto :parse_args
)
if "%~1"=="--build" (
    set "ACTION=build"
    shift
    goto :parse_args
)
if "%~1"=="-t" (
    set "ACTION=test"
    shift
    goto :parse_args
)
if "%~1"=="--test" (
    set "ACTION=test"
    shift
    goto :parse_args
)
if "%~1"=="-c" (
    set "ACTION=clean"
    shift
    goto :parse_args
)
if "%~1"=="--clean" (
    set "ACTION=clean"
    shift
    goto :parse_args
)
if "%~1"=="--docker" (
    set "USE_DOCKER=true"
    set "DEV_MODE=false"
    set "ACTION=start"
    shift
    goto :parse_args
)
if "%~1"=="--docker-dev" (
    set "USE_DOCKER=true"
    set "DEV_MODE=true"
    set "ACTION=start"
    shift
    goto :parse_args
)
if "%~1"=="--stop" (
    set "ACTION=stop"
    shift
    goto :parse_args
)
if "%~1"=="--restart" (
    set "ACTION=restart"
    shift
    goto :parse_args
)
if "%~1"=="--logs" (
    set "ACTION=logs"
    shift
    goto :parse_args
)
if "%~1"=="--health" (
    set "ACTION=health"
    shift
    goto :parse_args
)

call :log_error "未知选项: %~1"
call :show_help
exit /b 1

:help_and_exit
call :show_help
exit /b 0

REM 主逻辑
:main_logic
REM 如果没有指定动作，显示帮助
if "%ACTION%"=="" (
    call :show_help
    exit /b 0
)

REM 检查依赖
call :check_dependencies
if errorlevel 1 exit /b 1

REM 初始化环境
if not "%ACTION%"=="stop" if not "%ACTION%"=="logs" if not "%ACTION%"=="health" (
    call :init_environment
    if errorlevel 1 exit /b 1
)

REM 执行动作
if "%ACTION%"=="start" (
    if "%DEV_MODE%"=="true" (
        call :start_dev
    ) else (
        call :start_prod
    )
) else if "%ACTION%"=="build" (
    call :build_app
) else if "%ACTION%"=="test" (
    call :run_tests
) else if "%ACTION%"=="clean" (
    call :clean_build
) else if "%ACTION%"=="stop" (
    call :stop_services
) else if "%ACTION%"=="restart" (
    call :restart_services
) else if "%ACTION%"=="logs" (
    call :show_logs
) else if "%ACTION%"=="health" (
    call :health_check
)

goto :eof

REM 主程序入口
call :parse_args %*
