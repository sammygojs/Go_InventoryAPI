# One iota Technical Tasks

## Solution Architecture Task

A client has multiple systems that they want to integrate into their mobile app: an e-commerce platform, a shipment system, and a payment provider. Your job is to think about the best way to serve these to a mobile app via a single anonymous wrapper, using services common to cloud hosting providers. You can use which ever cloud provider you wish for this task.

Please draw up an integration diagram showing the flow, from how the mobile app will interact with your wrapper, to fetch and send data to these client systems.

## Product REST API Task

For this task, we would like you to create two endpoints for querying data. We have provided you with a JSON file containing test data.

### Endpoints
- For the first endpoint, you will need to return a list of products using the test data. We would also like you to show how many products there are in general in the response. You should also be able to reduce the amount of products coming back using a `limit` query parameter.
- For the second endpoint, you will need to return a single product based on an id path parameter.

#### Extra Tasks
- We would like you to handle a `locale` being provided in the request, which will change the descriptions and features of the product based on the given `locale`.
- We would also like to see you displaying membership pricing, if it is less than the standard, based on the user that is making the request. We should be able to make a request against both endpoints as a non-member and member user and see the correct pricing.
- In the get products endpoint we would like to be able to filter products based on: `minPrice`, `maxPrice`, `inStock`, `colour` (if a product has multiple colours, it should return for both, e.g. if the product colour is Red/Black then it will return when querying `colour=red` or `colour=black`).

All of the models that we would like you to map the data against have been provided. Please try to fill as much data as you can in these models. Along with this we have also provided you with a basic starting point for this task, which you can find in the `main.go` file.

Once you have completed the task, please submit your code, any tests, and instructions for running the the REST API task, as well as your architecture diagram.

If you have any questions feel free to get in touch with us via email at [recruitment@oneiota.co.uk](mailto:recruitment@oneiota.co.uk).

Good luck!

## FAQs

##### Which language should I use?
The role is mostly Go so you should try to undertake the task using this language.

##### Which cloud provider would be best for the architecture task?
We use AWS, so if you want to base the architecture task on something similar to what we would do, that would be great. But the decision is yours.

# ProductsAPI
