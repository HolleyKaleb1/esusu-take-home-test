#Esusu Memes as a Service Overview

Memes As A Service

Memes as a service API 
- Uses Json file as database for the sake of this project in the real world you have many options here you could deploy a db on a cloud service and use environment variables to gain access in your data access layer of your application
- Also uses util functions to implement the validations against longitude and latitude.
- For the sake of this application I decided to use dummy implementation of memes to save time however we could easily use an external API to get these memes in this case I have set a category as the query in the real world we may want to search metadata or allow ai to determine what is shown in memes and return them.This is up to the implementer of the service and what they are trying to achieve.
- Authorization
  - I am sending the user id as the bearer token. In the real world we would use jwts and get the users user id from the context of that request from the jwt that is sent with each request and check the db for their token balance if not there we will add 150 tokens for that user. We also support updating tokens forward and backwards from external service. External service could actually use a webhook to call our api sort of like stripe does with subscriptions for applications.
- 

API
- GET /memes
  - arguments 
    - long
      - optional
    - lat
      - optional
    - query
      - optional
  - response
  -  {
        "id": 1,
        "title": "Default Meme 1",
        "image": "https://example.com/meme1.jpg",
        "source": ""
    }

- POST /update-token-balance
  - arguments
    - {"user_id": "user1234", "amount": 50}
  - response
  - {
    "success": true,
    "message": "Token balance updated successfully",
    "user_id": "user1234",
    "new_balance": 150
    }
         

Answers to Questions 3-4
3a. 
  - In the case of this solution we would need to implement a robust CI/CD pipeline that automates the build testing and deploying processes. In order to ensure that new changes haven't impacted our performance we could include load testing in this process setting threshholds that fails pipelines if our API is not performing within our SLA. We could achieve this by using tools such as Azure pipelines along with Azure load testing. As well as unit testing and integration testing.

3b. 
  - How I would approach finding SLA's
      - Define clear SLAs for the API service, including response time objectives, availability targets, and error rate thresholds. (Load testing can achieve this)
      - Implement monitoring and alerting systems to track performance metrics and ensure compliance with SLAs. (Datadog,Splunk are tools that we can use to achieve this)
      - SLA's
        - Availability:
          - Target Availability: 99.9% or higher.
          - Downtime Allowance: Maximum of 4 minutes and 22 seconds of downtime per month.
        
        - Performance:
          - Average response time: 10 milliseconds (ms).
          - Maximum response time: 50 ms.
          - Throughput: Handle sustained traffic of 10,000 RPS with spikes up to 15,000 RPS.
          - Error Rate:
              - Maximum 0.1% of requests resulting in 5xx server errors.
              - Maximum 1% of requests resulting in client errors (e.g., 4xx errors).

3c. 
  - How I would support diverse clients and keep token tracking as fast as possible.
    - Deploy the  service across multiple geographical regions to reduce latency and improve the user experience for clients located in different parts of the world.
    - Use CDNs to cache static assets and reduce the load on servers.
    - Implement global load balancing and distribute traffic across distributed instances of the service.
    - Use distributed data storage solutions with built-in replication and consistency mechanisms to ensure that token balances are synchronized across regions without sacrificing performance.
    - Implement efficient caching strategies to minimize database read operations and reduce latency.
   
4. To modify the service to keep track of if a client/user has authorization to get AI generated memes we can use the following tools to achieve this 
  - Claims - This claim will be responsible of letting our services know what access a user has within our system (Stateless authentication)
    - An alternate solution to this is white listing a users id and using context of the request to get this data when api is called
    - *** Of course this is along with Authorization and authentication for a user / Service
  - Request Handling
    - With Claims being implemented we can now update the API's to use the claim to support different paths of our application (i.e. path for regular memes and path for ai memes).
  - Cache Implementation
    - Since we have now added complexity around the claims we can now implement a cache that caches/lazyloads the status of users, when user is updated we update the cache to ensure that performance is not impacted by claim changes
  - Also we will need to scale up our authorization system to scale up as our load increases

  Here is a potential design of the system at a high level
  - Flow
    - Client calls meme service
    - Load balancer handles the traffic based on the load that is needed based on load on service
    - Meme service takes token and uses authentication middleware to get context of jwt sent by client (Cache sits in front of this)
      - If the user is frequent or recent we will try cache if there is a cache miss we will go to the auth service and get the claims and hydrate the cache.
    -  Meme service will then check the tokens to see if user has enough tokens, will return that the user has insufficient tokens or if they have enough tokens we will then process the request and update the tokens after because we only want to remove tokens if the process was completed by our system. 
    -  Based on claims determine the path in application code whether you want to get default memes from db or if you want generated memes by AI.
      
  <img width="1534" alt="Screenshot 2024-05-05 at 3 46 35 PM" src="https://github.com/HolleyKaleb1/esusu-take-home-test/assets/42941354/888aed6f-1e96-41a6-8134-f195a017af24">

     

