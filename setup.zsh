function autoalias(){
    `pwd`/autoalias store $1
    `pwd`/autoalias render > ~/.autoaliases
    source ~/.autoaliases
}

autoload -Uz add-zsh-hook
add-zsh-hook preexec autoalias
