# FogFlow Roadmap
The FogFlow Generic Enabler is a distributed execution framework to support dynamic processing flows over cloud and edges.

## Introduction
This section elaborates on proposed new features or tasks which are expected to be added to the product in the foreseeable future.
There should be no assumption of a commitment to deliver these features on specific dates or in the order given. 
The development team will be doing their best to follow the proposed dates and priorities, but please bear in mind that plans to work on a given feature or task may be revised. 
All information is provided as general guidelines only, and this section may be revised to provide newer information at any time.

Disclaimer:
 1. This section has been last updated in October 2021. Please take into account its content could be obsolete.
 2. Note we develop this software in Agile way, so development plan is continuously under review. Thus, this roadmap has to be understood as rough plan of features to be done along time which is fully valid only at the time of writing it. This roadmap has not be understood as a commitment on features and/or dates.
 3. Some of the roadmap items may be implemented by external community developers, out of the scope of GE owners. Thus, the moment in which these features will be finalized cannot be assured.

### Short Term

The following list of features are planned to be addressed in the short term, and incorporated in a next release of the product:
1. Kubernetes support - 
   It will be usefull to use Fogflow in production environment. Following features would be implemented:   
   - High availability 
   - Self-Healing
   - Automated Rollouts & Rollback
   - K8s Security and Network Policy
   - Edge Node K8s Support
   - Worker thread implementation with K8s
   - Ease deployment with Helm
   - Taints and Tolerations
   - K8s Manual Creation

2. Use cases implementation of NGSI-LD  
   
3. User Manual Update
   - K8s Manual Creation
   - FogFlow user manual support for new feature and bugs. 

4. Improve Quality & Testing
   - Automation of Regression testing 
   - Performance Testing to evaluate the benchmarks
5. Edge AI
   - Support the edge node with Edge TPU.

### Medium Term
The following list of features are planned to be addressed in the medium term, typically within the subsequent release(s) generated in the next 6 months after the next planned release.
1. Firewall Support
   - Currently, the FogFlow edge node requires to have a public IP address to be accessible by the FogFlow cloud node. In the actual deployment environment, the FogFlow edge node is very often deployed behind the company firewall via NAT, we need to find a way to support this scenario. One way to address this is to find a proxy for such kind of edge nodes. For example, assign the FogFlow cloud broker to be the proxy for the FogFlow brokers at this type of edge nodes.
2. App repository and catalog. 
   - Create an app repository.
   - Add a wide range of data connectors and convertors
3.Performance Analysis
   - FogFlow REST APIs(Create, Update and subscribe) performance comparison with other GEs under different network environment.
4. Common Communication Channel
   - FogFlow control object should communicate with common channel. RabbitMQ should be only communication channel for ease and unicity.
   - Internal objects which user can define using the Designer GUI needs to be stored in Plain JSON format. It shold be protocol independent which is currently in NGSI-V1 format.
   - All the use cases need to be modified to full fill the protocol independent format.
5. Catalog creation for different operators 
    - FogFlow need to create different operators of NGSI-LD also like NGSI-V1 that any user can use according to their use case.  
    - The use case catalogue of operators for ready to play kind of use cases need to be added.
6. geoJson property location support
   - FogFlow need to provide support of geo Property in NGSI-LD, which will include the below support
     a) Creation and Update of NGSI-LD entity with geo Property.
     b) Provide support of geo based query.
     c)	Provide support of geo query in discovery.
  

### Long term
The following list of features are proposals regarding the longer-term evolution of the product even though the development of these features has not yet been scheduled for a release in the near future. Please feel free to contact us if you wish to get involved in the implementation or influence the roadmap:

1. Designer enhancement
   - FogFlow should add the support of control objects visualization. Need to provide monitoring support to visualize different devices from Designer web page. 
   - Need to add support for instant query from designer web page like how many devices are connected with which edge/cloud node.
2. Semantics-based data integration
   - creating dynamic data processing pipelines to convert arbitrary raw data into standard-based entities. 
3. Semantics-based service composition
   - linking serverless functions based on their semantically-annotated inputs and outputs.   
4. Digital twin support
   - make the current programming model to support the creation of digital twins and also the interaction between digital twins. 
5. Make FogFlow platform agnostic, Support of public clouds.
6. One Click deployment based on customer requirement.
7. Modularize the orchestration algorithms to achieve various service level objectives defined in a service intent.
8. Vocabulary management 
   - Manage all data types and attributes via Dgraph.
9. Multi-tenancy support.
   - Support multiple users over the same cloud-edge infrastructure.
