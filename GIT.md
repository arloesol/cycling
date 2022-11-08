# GIT

The repository is on Github: [cycling repo](https://github.com/arloesol/cycling)

## branches

- prod - deployed on https://cycling.arloesol.com
- preprod - deployed on https://cycling.beta.arloesol.com/
- main - main/central branch
- webedits - used to edit pages from internet
- feat-issuenbr-title - feature branches
- cont-issuenbr-title - content branches
- bf-issuenbr-title - bugfix branches

## GIT flow
### feature, bugfix and content branches

The main git flow used is based on [Github's flow](https://docs.github.com/en/get-started/quickstart/github-flow)

1. For any new changes (features, content, bugfix) a new branch is created from the main branch. 
1. Any changes to this new branch can be commited/pushed to the github repo. 
1. Once the changes are ready, a pull request into main is created after merging the main branch into the new branch. The pr mentions it closes the issue nbr : "closes #nbr" -> github will close the issue automatically once the pr is accepted
1. The pull request is accepted and squashed (ff only) into main
1. The original branch is deleted.  

### preprod and prod branches

The git flow for deployment is under review - currently following process is used

For deployment to preprod

1. pull request from main to preprod
1. accept pull request

For deployment to prod the following process is followed 

1. pull request from preprod to prod
1. accept pull request

### webedits branch

For the webedits branch a specific GIT flow still needs to be defined

## Deployment process

render.com checks commits on prod and preprod branches and starts a deploy automatically

- [production url](https://cycling.arloesol.com)
- [pre-production url](https://cycling.beta.arloesol.com)