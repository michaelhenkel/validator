go build
./validator print -n contrail > ~/mygraph.dot
dot -Tps  ~/mygraph.dot -o ~/mygraph.ps
dot -Tpng -Gdpi=300  ~/mygraph.dot > ~/g.png
open ~/g.png
