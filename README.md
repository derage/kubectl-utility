# kubectl-utility

## GOAL

I really love using [fubectl](https://github.com/kubermatic/fubectl). It makes interacting with a kubernetes cluster very easy. Being able to do shorthands makes for fast operations.

The down side with fubectl is it requires more then just kubectl to be installed, so getting this to work on something like a windows machince becomes hard. The aim for this project is to come up with a fubectl like utility, but instead have it all written in go, and only using libraries that can be crossed compiled.

### Should I include kubectl in this

So its written in golang, which means maybe I should just extend kubectl with the shorthands? I thought about this, but I dont want to have to build a new version every time a new version of kubectl is released. I want to try and make this agnostic to kubernetes version. So for now I will just have it be an external dependency