# kbn
### Learning Bayesian Networks with Bounded Treewidth and Latent Variables

___

## Using
To use **kbn** you need to install libgsl and some go packages described bellow.

* Install [libgsl](https://www.gnu.org/software/gsl/):

    In Ubuntu you can use the command bellow:

                apt install gsl-bin libgsl-dev

    For more details, check [gogsl](https://github.com/dtromb/gogsl) package.

* Get other required packages:

                go get github.com/dtromb/gogsl...
                go get github.com/willf/bitset...
                go get github.com/britojr/tcc...

* Get, test and install:

                go get github.com/britojr/kbn...
                go test github.com/britojr/kbn...
                go install github.com/britojr/kbn...

* Using the commands:

                learn --help

---
