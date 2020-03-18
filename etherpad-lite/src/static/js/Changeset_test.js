var Changeset = require("./Changeset");


var cs = "Z:z>1|2=m=b*0|1+1$\n";
var unpacked = Changeset.unpack(cs);

var opiterator = Changeset.opIterator(unpacked.ops);
console.log(opiterator);
while (opiterator.hasNext()){
    aOp = opiterator.next();
    console.log(aOp);
}

var AttribPool = require("./AttributePool");
var apool = new(AttribPool);
console.log(apool);
apool.fromJsonable({"numToAttrib":{"0":["author","a.kVnWeomPADAT2pn9"],"1":["bold","true"],"2":["italic","true"]},"nextNum":3});
console.log(apool);
console.log(apool.getAttrib(1));

var atext = {"text":"bold text\nitalic text\nnormal text\n\n","attribs":"*0*1+9*0|1+1*0*1*2+b|1+1*0+b|2+2"};
console.log(atext);
var opiterator = Changeset.opIterator(atext.attribs);
console.log(opiterator);
while (opiterator.hasNext()){
    console.log(opiterator.next());
}

