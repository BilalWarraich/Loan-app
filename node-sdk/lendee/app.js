'use strict';

//get libraries
const express = require('express');
var queue = require('express-queue');
const bodyParser = require('body-parser');
const request = require('request');
const path = require('path');

//create express web-app
const app = express();
const invoke = require('./invokeNetwork');
const query = require('./queryNetwork');
var _time = "T00:00:00Z";

//declare port
var port = process.env.PORT || 8000;
if (process.env.VCAP_APPLICATION) {
  port = process.env.PORT;
}

app.use(bodyParser.json({
  limit: '50mb', 
  extended: true

}));

app.use(bodyParser.urlencoded({
 limit: '50mb', 
 extended: true
}));

app.use(function(req, res, next) {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
  });

//Using queue middleware
app.use(queue({ activeLimit: 30, queuedLimit: -1 }));

//run app on port
app.listen(port, function () {
  console.log('app running on port: %d', port);
});

//-------------------------------------------------------------
//----------------------  POST API'S    -----------------------
//-------------------------------------------------------------

app.post('/api/addLendee', async function (req, res) {

  var request = {
    chaincodeId: 'loan',
    fcn: 'addLendee',
    args: [

      req.body.lendeeID,
      req.body.username,
      req.body.password

    ]
  };
console.log(req.body);
  let response = await invoke.invokeCreate(request);
  if (response) {
    if(response.status == 200)
    res.status(response.status).send({ message: "The Lendee with ID: "+req.body.lendeeID+ " is stored in the blockchain with " +response.message  });
    else
    res.status(response.status).send({ message: response.message});
  }
});

app.post('/api/addRequest', async function (req, res) {

  var request = {
    chaincodeId: 'loan',
    fcn: 'addRequest',
    args: [
      
      req.body.requestID,
      req.body.lendeeID,
      req.body.lenderID,
      req.body.loanRequired,
      req.body.product,
      req.body.returnDate,
      req.body.expectedOutcome,
      req.body.status,

    ]
  };
console.log(req.body);
  let response = await invoke.invokeCreate(request);
  if (response) {
    if(response.status == 200)
    res.status(response.status).send({ message: "The Loan Request with ID: "+req.body.requestID+ " is stored in the blockchain with " +response.message  });
    else
    res.status(response.status).send({ message: response.message});
  }
});

app.post('/api/updateResponse', async function (req, res) {

  var request = {
    chaincodeId: 'loan',
    fcn: 'updateResponse',
    args: [

      req.body.responseID,
      req.body.newStatus

    ]
  };
console.log(req.body);
  let response = await invoke.invokeCreate(request);
  if (response) {
    if(response.status == 200)
    res.status(response.status).send({ message: "The Lendee Response to Lender with ID: "+req.body.responseID+ " is stored in the blockchain with " +response.message  });
    else
    res.status(response.status).send({ message: response.message});
  }
});


//-------------------------------------------------------------
//----------------------  GET API'S  --------------------------
//-------------------------------------------------------------

app.get('/api/queryLendee', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryLendee',
    args: [
      req.query.username,
      req.query.password
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
      res.status(response.status).send(JSON.parse(response.message));
    else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryResponse', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryResponse',
    args: [
      req.query.status,
      req.query.lendeeID
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
      res.status(response.status).send(JSON.parse(response.message));
    else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryResponses', async function (req, res) {
  var data = [];

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryResponse',
    args: [
      req.query.status,
      req.query.lendeeID
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200){
      
      var resp = JSON.parse(response.message)
      console.log(resp);
      for (var i=0 ; i<resp.length ; i++){
      var lenderID = resp[i].Record.lenderID;
      const request1 = {
        chaincodeId: 'loan',
        fcn: 'queryLenderByID',
        args: [
             lenderID
            
        ]
      };
      console.log(req.query);
      let response1 = await query.invokeQuery(request1)
      if (response1) {
        if(response1.status == 200){
        var resp1 = JSON.parse(response1.message)
        console.log(resp[0].Record);
        var obj = Object.assign(resp[i].Record,resp1[0].Record);
        console.log(resp);
        console.log(obj);
        data.push(obj);
      }
    }
    }
    res.status(response.status).send(data);

    }
      else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryRequests', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryRequests',
    args: [
      req.query.lendeeID,
      req.query.lenderID
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
      res.status(response.status).send(JSON.parse(response.message));
    else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryAllLenders', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryAllLenders',
    args:[]
   
  };
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
    res.status(response.status).send(JSON.parse(response.message));
    else
    res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryLenderByName', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryLenderByName',
    args: [
      req.query.username
        
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
      res.status(response.status).send(JSON.parse(response.message));
    else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryRequestByID', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryRequestByID',
    args: [
      req.query.lendeeID
        
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
      res.status(response.status).send(JSON.parse(response.message));
    else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryLenderByID', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryLenderByID',
    args: [
      req.query.lenderID
        
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
      res.status(response.status).send(JSON.parse(response.message));
    else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryResponseID', async function (req, res) {

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryResponseID',
    args: [
      req.query.responseID
        
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200)
      res.status(response.status).send(JSON.parse(response.message));
    else
      res.status(response.status).send({ message: response.message });
  }
});

app.get('/api/queryExistingRequests', async function (req, res) {
  var data = [];

  const request = {
    chaincodeId: 'loan',
    fcn: 'queryRequestByID',
    args: [
      
      req.query.lendeeID
    ]
  };
  console.log(req.query);
  let response = await query.invokeQuery(request)
  if (response) {
    if(response.status == 200){
      
      var resp = JSON.parse(response.message)
      console.log(resp);
      for (var i=0 ; i<resp.length ; i++){
      var lenderID = resp[i].Record.lenderID;
      const request1 = {
        chaincodeId: 'loan',
        fcn: 'queryLenderByID',
        args: [
             lenderID
            
        ]
      };
      console.log(req.query);
      let response1 = await query.invokeQuery(request1)
      if (response1) {
        if(response1.status == 200){
        var resp1 = JSON.parse(response1.message)
        console.log(resp[0].Record);
        var obj = Object.assign(resp[i].Record,resp1[0].Record);
        console.log(resp);
        console.log(obj);
        data.push(obj);
      }
    }
    }
    res.status(response.status).send(data);

    }
      else
      res.status(response.status).send({ message: response.message });
  }
});