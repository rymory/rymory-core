'''
// Copyright (c) 2017-2026 Onur Yaşar
// Licensed under AGPL v3 + Commercial Exception
// See LICENSE.txt
'''

from http import HTTPStatus
import os
from sendgrid import SendGridAPIClient # type: ignore
from sendgrid.helpers.mail import Mail # type: ignore
import http.client
import json

def main(args):
    '''
    Takes in the email address, subject, and message to send an email using SendGrid, 
    returns a json response letting the user know if the email sent or failed to send.

        Parameters:
            args: Contains the from email address, to email address, subject and message to send

        Returns:
            json body: Json response if the email sent Successfully or if an error happened
    '''

    print("burda")
   
    valid_api_url  = os.getenv('VALID_API_URL')

    connection = http.client.HTTPSConnection('api.github.com')

    headers = {'Content-type': 'application/json'}

    foo = {'text': 'Hello world github/linguist#1 **cool**, and #1!'}
    json_foo = json.dumps(foo)

    connection.request('GET', valid_api_url, None, headers)

    response = connection.getresponse()

    user_id = response.getheader("userid")

    if user_id != None:
        return {
            "statusCode" : HTTPStatus.BAD_REQUEST,
            "body" : "no valid user authontication"
        }
    
    key = os.getenv('API_KEY')
   
    user_from = args.get("from")
    user_to = args.get("to")
    user_subject = args.get("subject")
    content = args.get("content")

    if not user_from:
        return {
            "statusCode" : HTTPStatus.BAD_REQUEST,
            "body" : "no user email provided"
        }
    if not user_to:
        return {
            "statusCode" : HTTPStatus.BAD_REQUEST,
            "body" : "no receiver email provided"
        }
    if not user_subject:
        return {
            "statusCode" : HTTPStatus.BAD_REQUEST,
            "body" : "no subject provided"
        }
    if not content:
        return {
            "statusCode" : HTTPStatus.BAD_REQUEST,
            "body" : "no content provided"
        }

    sg = SendGridAPIClient(key)
    message = Mail(
        from_email = user_from,
        to_emails = user_to,
        subject = user_subject,
        html_content = content)
    response = sg.send(message)

    if response.status_code != 202:
        return {
            "statusCode" : response.status_code,
            "body" : "email failed to send"
        }
    return {
        "statusCode" : HTTPStatus.ACCEPTED,
        "body" : "0x11031:Success"
    }

args = {'from': 'Hello world github/linguist#1 **cool**, and #1!'}
main(args)