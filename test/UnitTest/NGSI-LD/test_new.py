import os,sys
# change the path accoring to the test folder in system
from datetime import datetime
import copy
import json
import requests
import time
import pytest
import ld_data
import sys

# change it by broker ip and port
brokerIp="http://localhost:8070"
discoveryIp="http://localhost:8090"

# test if header content-Type application/json is allowed or not 
def test_case74():
	url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
	headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
	r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
	assert r.status_code == 204


# test if header content-Type is application/ld+json then the link header should not be persent in request

def test_case75():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/ld+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
        assert r.status_code == 404

#test if Allowd Content-Type are only appliation/json and application/ld+json

def test_case76():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application1/ld1+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
	print(r.status_code)
        assert r.status_code == 400


# test create and get the entity in openiot FiwareService

def test_case77():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
        #print(r.status_code)
	url=brokerIp+'/ngsi-ld/v1/entities/'+'urn:ngsi-ld:Device:water001'
	r = requests.get(url,headers=headers)
	assert r.status_code == 200


def test_case78():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application1/ld1+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.testData74),headers=headers)
        #print(r.status_code)
	headers={'Content-Type' : 'application1/ld1+json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiott','fiware-servicepath' :'test'}

        url=brokerIp+'/ngsi-ld/v1/entities/'+'urn:ngsi-ld:Device:water001'
        r = requests.get(url,headers=headers)
        assert r.status_code == 404

# To test upsert Api support only array of entities
def test_case79():
        url=brokerIp+"/ngsi-ld/v1/entityOperations/upsert"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"', 'fiware-service' : 'openiot','fiware-servicepath' :'test'}
        r=requests.post(url,data=json.dumps(ld_data.testData75),headers=headers)
        #print(r.status_code)
	assert r.status_code == 500


def test_case48():
        #to create an entity
        url=brokerIp+"/ngsi-ld/v1/entities/"
        headers={'Content-Type' : 'application/json','Accept':'application/ld+json','Link':'<{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'}
        r=requests.post(url,data=json.dumps(ld_data.subdata38),headers=headers)
        #print(r.content)
        #print(r.status_code)

        #to fetch the registration from discovery
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A3000"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp["AttributesList"]["https://uri.etsi.org/ngsi-ld/default-context/brandName"])
        print("\nchecking if brandName attribute is present in discovery before deletion")
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A3000":
                if resp["AttributesList"]["https://uri.etsi.org/ngsi-ld/default-context/brandName"]["type"] == "Property":
                        print("\n-----> brandName is existing...!!")
                else:
                        print("\n-----> brandName does not exist..!")
        else:
                print("\nNot Validated")
        #print(r.status_code)


        #to delete brandName attribute
        url=brokerIp+"/ngsi-ld/v1/entities/urn:ngsi-ld:Vehicle:A3000/attrs/brandName"
        r=requests.delete(url)
        #print(r.content)
        #print(r.status_code)

        #To fetch registration again from discovery
        url=discoveryIp+"/ngsi9/registration/urn:ngsi-ld:Vehicle:A3000"
        r=requests.get(url)
        resp_content=r.content
        resInJson= resp_content.decode('utf8').replace("'", '"')
        resp=json.loads(resInJson)
        #print(resp["AttributesList"])
        print("\nchecking if brandName attribute is present in discovery after deletion")
        if resp["ID"]=="urn:ngsi-ld:Vehicle:A3000":
                if "https://uri.etsi.org/ngsi-ld/default-context/brandName" in resp["AttributesList"]:
                        print("\n-----> brandName is existing...!!")
                else:
                        print("\n-----> brandName does not exist because deleted...!")
        else:
                print("\nNot Validated")
        assert r.status_code == 200


