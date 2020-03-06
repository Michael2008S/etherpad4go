// index.js
const express = require('express');
const path = require('path');
const app = express();
// 在 app 文件夹开启静态服务
app.use(express.static('app'));
app.listen(8090, () => {
  console.log('Demo server listening on port 8090');
});