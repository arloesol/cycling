# Tools

## Website scraping

The content of the routes comes from other websites

This content is scraped by using go code in /tools/scrape

for example the www.visitlimburg.be website can be scraped by running the code found in /tools/scrape/limburg

```shell
cd tools/scrape/limburg
go run limburg.go
```

The resulting content, gpx files and images will be in the route, gpx and img directories

They can be deployed to the right directories using following command

```bash
cd bin
./deploy.sh limburg
```

## Git flow 

The git flow used in the repo is described [here](GIT.md)

Following bash scripts can be used for several steps in the development process

- new development : bin/git_newdev.sh *issuenr*
- close a development : bin/git_closedev.sh
- switch to other branch : bin/git_switch.sh
- clean unused branches : bin/git_cleanup.sh (still in beta ... beware)

some of these scripts use the github cli - install with apt and then login to github 

```bash
sudo apt -y install gh
gh auth login
```

## Other

/bin/install.sh can be used to install some useful items related to this project

## .bashrc

Add the following to your ~/.bashrc file

change the "gitdirectory" value to your local dir

```bash
export GITDIR=$HOME/gitdirectory/cycling
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GITDIR/bin

eval "$(gh completion -s bash)"
```
