### Setup

To run the project you will need:

- Install Go v1.13.*
- Install MySQL 8.*
- Export your `$GOBIN` to `$PATH` in `.bash_profile | .zshrc`: `export PATH="$PATH:$GOBIN"`
- Setup ENV variables from `.env.dist`, you should automate that with [`direnv`](https://direnv.net/)

The **Makefile** provides all the useful commands to run and test the project.
