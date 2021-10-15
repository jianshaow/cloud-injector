
function blue() {
    echo -e "\e[34m$1\e[0m"
}

function green() {
    echo -e "\e[32m$1\e[0m"
}

function red() {
    echo -e "\e[31m$1\e[0m"
}

function h_blue() {
    echo -e "\e[34m\e[01m$1\e[0m"
}

function h_green() {
    echo -e "\e[32m\e[01m$1\e[0m"
}

function h_red() {
    echo -e "\e[31m\e[01m$1\e[0m"
}

function h_yellow() {
    echo -e "\e[33m\e[01m$1\e[0m"
}

function b_red() {
    echo -e "\e[31m\e[01m\e[05m$1\e[0m"
}

function b_yellow() {
    echo -e "\e[33m\e[01m\e[05m$1\e[0m"
}
