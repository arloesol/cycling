# GIT

The repository is on Github: [cycling repo](https://github.com/arloesol/cycling)

## Branches

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

Following scripts are used for this process [git_newdev.sh](/bin/git_newdev.sh) and [git_closedev.sh](/bin/git_closedev.sh)

### preprod and prod branches

The git flow for deployment is still under review - currently the following process is used

For deployment to preprod

1. align preprod with main branch
1. the release is tagged with a version nbr of type vn.m.o (with o increased by 1)
1. bash script [git_main_to_preprod.sh](/bin/git_main_to_preprod.sh) is used for this

For deployment to prod the following process is followed

1. align prod with preprod branch
1. the release is tagged with a version nbr of type vn.m (with m increased by 1)
1. bash script [git_preprod_to_prod.sh](/bin/git_preprod_to_prod.sh) is used for this

### webedits branch

For the webedits branch a specific GIT and process flow still needs to be defined

Some changes and commits will be made via github gui on this branch.

## Tags

releases to preprod and prod are tagged with version numbers

each release to preprod gets a version number of type vm.n.o with o increased at each release in preprod - f.ex. v1.2.3, v1.2.4, v1.2.5 ...

each release to prod gets a version number of type vm.n with n increased at each release in prod - f.ex. v1.2, v1.3, v1.4 ...

once a release is done to prod (always from preprod) - the preprod version gets aligned with the prod version - f.ex. v1.2.3 -> v1.2.4 (preprod release) -> v1.3 (prod release) -> v1.3.1 (preprod release) -- v1.2.4 and v1.3 are the same version

increases in major release number will be handled manually by tagging the specific commit
## Deployment process

render.com checks commits on prod and preprod branches and starts a deployment automatically

The status of the deployments can be seen [here](https://dashboard.render.com/)
