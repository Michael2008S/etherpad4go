var Changeset = require("./Changeset");


var cs = "Z:z>1|2=m=b*0|1+1$\n";
var unpacked = Changeset.unpack(cs);

var opiterator = Changeset.opIterator(unpacked.ops);
console.log(opiterator);
var aOp = opiterator.next();
console.log(aOp);
aOp = opiterator.next();
console.log(aOp);
hasNext = opiterator.hasNext(aOp);
console.log(hasNext);
aOp = opiterator.next();
console.log(aOp);
aOp = opiterator.next();
console.log(aOp);

hasNext = opiterator.hasNext(aOp);
console.log(hasNext);

aOp = opiterator.next();
console.log(aOp);


