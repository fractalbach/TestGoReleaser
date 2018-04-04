#!/bin/bash

# -----------------------------------------------------------------------
# USAGE:
#                   manybuild.sh [Project Name]
#
# Customizations:
# 
#
#   DIR:            The target directory where the builds will go.
#   PROJECT:        The name of the file, and its zipped container.
#       
#   AMD64_LIST:     Names of operating systems for amd64 architecture.
#   ARM_LIST:       All Linux, the versions of arm architecture.
#
#
# 


PROJECT=${1}
DIR="releases"
AMD64_LIST="windows linux darwin"
ARM_LIST="6 7"

# -----------------------------------------------------------------------
# Function Definitions
# ________________________________________


function getName {
    local OS=${1}
    local ARCH=${2}
    local ARM=${3}
    local EXT=${4}
    echo ${DIR}/${PROJECT}_${OS}_${ARCH}${ARM}${EXT}
}



function buildMe {
    local OS=${1}
    local ARCH=${2}
    local ARM=${3}
    local EXT=""
    local OutputPath=""

    # Windows requires the file extension ".exe" at the end of the file-path.
    #
    if [[ "$OS" = "windows" ]]; then
        local EXT=".exe"
    fi


    # Create a unique and descriptive file-path for the binary file.
    
    local OutputPath=${DIR}/${PROJECT}_${OS}_${ARCH}${ARM}${EXT}


    # Check for any missing arguments, or anything that might prevent 
    # the program from building correctly.
    
    if [[ "$OutputPath" = "" ]]; then
        echo "Error: Empty Output Path."
        return
    fi
    

    # Run the "go build" command under the customized environment.
    
    env GOOS=${OS} ARCH=${ARCH} GOARM=${ARM} go build -v -o ${OutputPath}

}


# ------------------------------------------------------------------
# Main Process. 
#
function RunMyCustomBuild {

    echo "Starting Customized Compilation."
    echo "Directory: "${DIR}
    echo "Project:   "${PROJECT}
    echo ""

    for i in $AMD64_LIST; do
        echo "Building:  OS: ("${i}"),  ARCH: (amd64)"
        buildMe ${i} "amd64"
    done

    for i in $ARM_LIST; do
        echo "Building:  OS: (Linux),  ARCH: (arm),  ARM: ("${i}")" 
        buildMe "linux" "arm" ${i}
    done

    echo "Custom Compile Script has Finished."
}

RunMyCustomBuild
# ------------------------------------------------------------------

