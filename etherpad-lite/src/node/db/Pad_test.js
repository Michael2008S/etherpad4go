var Pad = require("./Pad");
var revNum = 2;
console.log(Pad);
var apad = new Pad.Pad("testpad");
console.log(apad);
var res = apad.getKeyRevisionNumber(revNum);
console.log(res, revNum);

console.log(Math.floor(51 / 100) * 100);


var Changeset = require("ep_etherpad-lite/static/js/Changeset");

text = `Welcome to Etherpad!\n\nThis pad text is synchronized~ https:\/\/github.com\/ether\/etherpad-lite\n`;
// text = "Welcome to Etherpad!\\n\\nThis pad text is synchronized~ https:\\/\\/github.com\\/ether\\/etherpad-lite\\n";
let cleanTxt = Pad.cleanText(text);
console.log(text);
console.log("cleanTxt:", cleanTxt);
let firstChangeset = Changeset.makeSplice("\n", 0, 0, cleanTxt);
console.log(firstChangeset);
atext = Changeset.makeAText("");
console.log(atext, apad.pool);
newAText = Changeset.applyToAText(firstChangeset, apad.atext, apad.pool);
console.log(newAText);

console.log("====================test new applyToAText==================")

reqCs = "Z:2l>1|3=2k*0+1$a";
reqAtext = Changeset.applyToAText(reqCs,newAText,apad.pool);
console.log(reqAtext);
// |3+2k*0+1|1+1
