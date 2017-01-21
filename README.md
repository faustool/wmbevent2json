# wmbevent2json
Go library to transform IBM WebSphere Message Broker/IntegrationBus monitoring XML messages to Json documents
The reason behind this project is to support the creation of an Elastic beat

##Example
The following wmb:event XML usually published by IBM Integration Bus:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<wmb:event xmlns:wmb="http://www.ibm.com/xmlns/prod/websphere/messagebroker/6.1.0/monitoring/event">
    <wmb:eventPointData>
        <wmb:eventData wmb:eventSchemaVersion="" wmb:eventSourceAddress="QueueName.terminal.Out"
                       wmb:productVersion="10.0.0.6">
            <wmb:eventIdentity wmb:eventName="Test Event" wmb:priority="" wmb:severity="" wmb:successDisposition=""/>
            <wmb:eventSequence wmb:counter="1" wmb:creationTime="2001-12-31T12:00:00" wmb:dataType="dateTime"
                               wmb:name="" wmb:value=""/>
            <wmb:eventCorrelation wmb:globalTransactionId="123" wmb:localTransactionId="456"
                                  wmb:parentTransactionId="789"/>
        </wmb:eventData>
        <wmb:messageFlowData>
            <wmb:broker wmb:UUID="" wmb:hostName="" wmb:name="MYBROKER"/>
            <wmb:executionGroup wmb:UUID="" wmb:name="MYEG"/>
            <wmb:messageFlow wmb:UUID="" wmb:name="MyFlowName" wmb:threadId="" wmb:uniqueFlowName="my.flow.name"/>
            <wmb:node wmb:detail="QUEUE.NAME" wmb:nodeLabel="QueueName" wmb:nodeType="ComIbmMQInputNode"
                      wmb:terminal="Out"/>
        </wmb:messageFlowData>
    </wmb:eventPointData>
    <wmb:applicationData>
        <wmb:simpleContent wmb:dataType="string" wmb:name="simple1" wmb:targetNamespace="" wmb:value="value1"/>
        <wmb:complexContent wmb:elementName="Complex1" wmb:targetNamespace="">
            <Complex1>
                <Child1 attr="attr-value">
                    Child value
                </Child1>
            </Complex1>
        </wmb:complexContent>
        <wmb:simpleContent wmb:dataType="string" wmb:name="simple2" wmb:targetNamespace="" wmb:value="value2"/>
        <wmb:complexContent wmb:elementName="Complex2" wmb:targetNamespace="">
            <Complex2>
                <Child1>
                </Child1>
            </Complex2>
        </wmb:complexContent>
    </wmb:applicationData>
    <wmb:bitstreamData>
        <wmb:bitstream wmb:encoding="CDATA"><![CDATA[<greeting>Hello, world!</greeting>]]> </wmb:bitstream>
    </wmb:bitstreamData>
</wmb:event>
```

Generates the Json below to be sent out to an Elasticsearch instance:
```json
{
  "event_wmb":"http://www.ibm.com/xmlns/prod/websphere/messagebroker/6.1.0/monitoring/event",
  "eventData_eventSourceAddress":"QueueName.terminal.Out",
  "eventData_productVersion":"10.0.0.6",
  "eventIdentity_eventName":"Test Event",
  "eventSequence_counter":"1",
  "eventSequence_creationTime":"2001-12-31T12:00:00",
  "eventSequence_dataType":"dateTime",
  "eventCorrelation_globalTransactionId":"123",
  "eventCorrelation_localTransactionId":"456",
  "eventCorrelation_parentTransactionId":"789",
  "broker_name":"MYBROKER",
  "executionGroup_name":"MYEG",
  "messageFlow_name":"MyFlowName",
  "messageFlow_uniqueFlowName":"my.flow.name",
  "node_detail":"QUEUE.NAME",
  "node_nodeLabel":"QueueName",
  "node_nodeType":"ComIbmMQInputNode",
  "node_terminal":"Out",
  "bitstream_encoding":"CDATA",
  "bitstream":"<greeting>Hello,world!</greeting>",
  "simpleContent":[
    {
      "dataType":"string",
      "name":"simple1",
      "value":"value1"
    },
    {
      "dataType":"string",
      "name":"simple2",
      "value":"value2"
    }
  ],
  "complexContent":[
    {
      "elementName":"Complex1",
      "data":{
        "{}:Complex1":{
          "{}:Child1":{
            "@attr":"attr-value",
            "#text":"Childvalue"
          }
        }
      }
    },
    {
      "elementName":"Complex2",
      "data":{
        "{}:Complex2":{
          "{}:Child1":{

          }
        }
      }
    }
  ]
}
```
