var Changeset = require("../src/static/js/Changeset");
var AttributePool = require("../src/static/js/AttributePool");

var cs = "|1+l*1+2|2+1z*2+2|1+1";
var newPool = new AttributePool();

var oldPool = new AttributePool();
oldPool.numToAttrib = {
    '0': ['author', 'a.glqITynU8VYvF40s'],
    '1': ['author', 'a.UaSfrktmubohgvYq'],
};
oldPool.attribToNum = { 'author,a.glqITynU8VYvF40s': 0, 'author,a.UaSfrktmubohgvYq': 1 };
oldPool.nextNum = 2;

var newcs = Changeset.moveOpsToNewPool(cs,oldPool,newPool);

console.log(newcs,newPool);
