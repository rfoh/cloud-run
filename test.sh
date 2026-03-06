#!/bin/bash

# Script para executar testes automatizados com Docker/Docker-Compose
# Demonstra o funcionamento da aplicação

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║   Weather Application - Automated Test Suite               ║"
echo "║   Using Docker / Docker-Compose                            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${BLUE}→${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Opção 1: Testes Unitários
run_unit_tests() {
    print_step "Executando testes unitários..."
    docker-compose run --rm test-unit
    print_success "Testes unitários completados!"
}

# Opção 2: Testes E2E
run_e2e_tests() {
    print_step "Executando testes end to end..."
    docker-compose run --rm test-e2e
    print_success "Testes end to end completados!"
}

# Opção 3: Cobertura de Testes
run_coverage_tests() {
    print_step "Executando análise de cobertura de testes..."
    docker-compose run --rm test-coverage
    print_success "Cobertura de testes completada!"
}

# Opção 4: Build da Aplicação
build_app() {
    print_step "Build da aplicação..."
    docker-compose run --rm build-app
    print_success "Build completado!"
}

# Opção 5: Executar Aplicação
run_app() {
    print_step "Iniciando a aplicação..."
    echo ""
    echo -e "${YELLOW}A aplicação estará disponível em: http://localhost:8080${NC}"
    echo -e "${YELLOW}Exemplo de requisição:${NC}"
    echo -e "${YELLOW}  curl 'http://localhost:8080/?cep=01310100'${NC}"
    echo ""
    docker-compose up app
}

# Opção 6: Todos os testes
run_all_tests() {
    print_step "Executando todos os testes..."
    echo ""
    
    run_unit_tests
    echo ""
    
    run_e2e_tests
    echo ""
    
    run_coverage_tests
    echo ""
    
    print_success "Todos os testes foram executados com sucesso!"
}

# Opção 7: Limpeza
cleanup() {
    print_step "Limpando containers e imagens..."
    docker-compose down -v
    print_success "Limpeza completada!"
}

# Menu principal
show_menu() {
    echo ""
    echo "Selecione uma opção:"
    echo "  1) Testes Unitários"
    echo "  2) Testes E2E"
    echo "  3) Cobertura de Testes"
    echo "  4) Build da Aplicação"
    echo "  5) Executar Aplicação"
    echo "  6) Executar Todos os Testes"
    echo "  7) Limpeza (remover containers)"
    echo "  0) Sair"
    echo ""
}

# Se nenhum argumento for passado, mostrar menu interativo
if [ $# -eq 0 ]; then
    while true; do
        show_menu
        read -p "Opção: " choice
        case $choice in
            1) run_unit_tests ;;
            2) run_e2e_tests ;;
            3) run_coverage_tests ;;
            4) build_app ;;
            5) run_app ;;
            6) run_all_tests ;;
            7) cleanup ;;
            0) 
                echo "Saindo..."
                exit 0
                ;;
            *) 
                print_error "Opção inválida!"
                ;;
        esac
    done
else
    # Se um argumento for passado, executar correspondente
    case "$1" in
        unit)
            run_unit_tests
            ;;
        e2e)
            run_e2e_tests
            ;;
        coverage)
            run_coverage_tests
            ;;
        build)
            build_app
            ;;
        run)
            run_app
            ;;
        all)
            run_all_tests
            ;;
        clean)
            cleanup
            ;;
        *)
            echo "Uso: $0 {unit|e2e|coverage|build|run|all|clean}"
            echo ""
            echo "Exemplos:"
            echo "  $0 unit        - Executar testes unitários"
            echo "  $0 e2e         - Executar testes end to end"
            echo "  $0 coverage    - Executar análise de cobertura"
            echo "  $0 all         - Executar todos os testes"
            echo "  $0 run         - Executar a aplicação"
            echo "  $0 clean       - Limpeza de containers"
            exit 1
            ;;
    esac
fi
