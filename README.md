# kbn
### Learning Bayesian Networks with Bounded Treewidth and Latent Variables

___

## Installation and usage
To use **kbn** you need to install libgsl and some go packages described bellow.

* Install [libgsl](https://www.gnu.org/software/gsl/):

    In Ubuntu you can use the command bellow:

                apt install gsl-bin libgsl-dev

    For more details, check the [gogsl packgage](https://github.com/dtromb/gogsl) page.

* Get other required packages:

                go get github.com/dtromb/gogsl...
                go get github.com/willf/bitset...
                go get github.com/britojr/tcc...

* Get, test and install:

                go get github.com/britojr/kbn...
                go test github.com/britojr/kbn...
                go install github.com/britojr/kbn...

* Commands:

                learn --help
