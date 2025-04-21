#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

CD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$CD_DIR"

echo -e "${YELLOW}Starting local GitHub Actions runner...${NC}"

# Function to run a specific job
run_job() {
    job=$1
    echo -e "\n${YELLOW}Running job: ${GREEN}$job${NC}"
    
    if docker-compose up --build --force-recreate "$job"; then
        echo -e "\n${GREEN}✓ Job '$job' completed successfully${NC}"
        return 0
    else
        echo -e "\n${RED}✗ Job '$job' failed${NC}"
        return 1
    fi
}

# Function to run all jobs
run_all_jobs() {
    local failed=()
    local jobs=("lint" "test" "build" "docker-build")
    
    for job in "${jobs[@]}"; do
        if ! run_job "$job"; then
            failed+=("$job")
        fi
    done
    
    # Print summary
    echo -e "\n${YELLOW}========== GitHub Actions Simulation Summary ==========${NC}"
    
    if [ ${#failed[@]} -eq 0 ]; then
        echo -e "${GREEN}All jobs completed successfully!${NC}"
        return 0
    else
        echo -e "${RED}The following jobs failed:${NC}"
        for job in "${failed[@]}"; do
            echo -e "${RED}- $job${NC}"
        done
        return 1
    fi
}

# Show help
show_help() {
    echo -e "Usage: ./run-actions.sh [OPTION]"
    echo -e "Run GitHub Actions workflows locally."
    echo -e "\nOptions:"
    echo -e "  lint          Run only the lint job"
    echo -e "  test          Run only the test job" 
    echo -e "  build         Run only the build job"
    echo -e "  docker-build  Run only the Docker build job"
    echo -e "  all           Run all jobs (default if no option is provided)"
    echo -e "  clean         Clean up Docker resources"
    echo -e "  help          Show this help message"
}

# Clean up
clean_up() {
    echo -e "${YELLOW}Cleaning up Docker resources...${NC}"
    docker-compose down --volumes --remove-orphans
    docker-compose rm -f
    echo -e "${GREEN}Cleanup completed${NC}"
}

# Main execution
case "$1" in
    lint|test|build|docker-build)
        run_job "$1"
        ;;
    all|"")
        run_all_jobs
        ;;
    clean)
        clean_up
        ;;
    help)
        show_help
        ;;
    *)
        echo -e "${RED}Unknown option: $1${NC}"
        show_help
        exit 1
        ;;
esac 