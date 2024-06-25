sudo apt update -y && sudo apt upgrade -y
sudo chsh -s /bin/zsh vscode
curl -fsSL https://get.pnpm.io/install.sh | sh -
source /home/vscode/.zshrc
sudo pnpm install -g aws-cdk

# cdk --version
# cdk init app --language=go

# GOOS=linux GOARCH=amd64 go build -o bootstrap
# zip lambdafucntion.zip bootstrap

# cdk diff
# cdk deploy
