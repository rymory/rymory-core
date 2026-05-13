// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt

const AWS = require('aws-sdk');
const SPACES_ENDPOINT = process.env['SPACES_ENDPOINT']; 
const SPACES_NAME = process.env['SPACES_NAME'];
const SPACES_KEY = process.env['SPACES_KEY']
const SPACES_SECRET = process.env['SPACES_SECRET']
const VALID_API_URL = process.env['VALID_API_URL']

const s3 = new AWS.S3({
  endpoint: SPACES_ENDPOINT,
  accessKeyId: SPACES_KEY,
  secretAccessKey: SPACES_SECRET
})

var request = require('request');

const options = {
  method: 'GET',
  url:VALID_API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Authorization' : ''
  },
};

async function main(args) {
  if (args.__ow_method == "options") {
    return {
      headers: {
        'Access-Control-Allow-Methods': 'OPTIONS, POST',
        'Access-Control-Allow-Origin': '*'
      },
      statusCode: 200
    }
  }

options.headers.Authorization = args.http.headers.authorization;



const userId = await new Promise((resolve, reject) => {
  request(options, function(err, res, body) {
    
    resolve(res.headers['userid']);

  });
});


console.log(userId);

  if (userId == undefined)  {
    return {
      statusCode: 400,
      body: {
        message: '0x11032:User Id is empty',
      }
    }
  };

  
  const fileName = userId + "---" +args.file_name;
  const contentType = args.content_type;
  const oldFileName = args.old_file_name;

  if (!fileName || !contentType) {
    console.log('0x11033:missing file_name or content_type')
    return {
      statusCode: 400,
      body: {
        message: '0x11033:missing file_name or content_type',
      }
    }
  }

try {
  var oldKey = userId + oldFileName.split(userId)[1];
  s3.deleteObject({
    Bucket: SPACES_NAME,
    Key: oldKey
  },function (err,data){})
} catch {}

  const params = {
    Bucket: SPACES_NAME,
    Fields: {
      'Content-Type': contentType,
      key: fileName,
    },
    Expires: 300,
    Conditions: [
      { 'acl': 'public-read' }
    ]
  };

  try {
    const signedPayload = await new Promise((resolve, reject) => {
      s3.createPresignedPost(params, (err, data) => {
        if (err) {
          reject(err);
          console.error(data)
          return;
        }
        resolve(data);
      })
    })
    console.log(`0x11031:Successfully signed payload for ${signedPayload.url}/${signedPayload.fields.key}`)

    return {
      statusCode: 200,
      body: {
        payload: signedPayload,
      }
    }
  } catch (error) {
    return {
      statusCode: 400,
      body: {
        message: `0x11068:unable to get signed url: ${error.message}`,
      }
    }
  }

}

exports.main = main;

