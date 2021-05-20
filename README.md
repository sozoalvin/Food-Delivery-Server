<link rel="stylesheet" href="styles.css">

# Food Delivery REST API Backend Server
Full REST API Backend Server System for Food Delivery Applications. Written in go/ golang with search, auto complete and suggestive search as well. The server also utilizes different flags that can be switced on and off, depending on environment settings and testing requirements.

<img src = "https://i.imgur.com/XFUVxvV.png" width="700">

<h2>Introduction</h2>

<p>The KV app was created in an attempt to solve teething customer service issues that have been affecting customer’s average order values, customer satisfaction and ultimate customer loyalty.
The app used to be running on local memory but has since moved to AWS's RDS with the server powered by EC2, running on Ubuntu.</p>

<h2>Tech Stacks Employed</h2>
<p>Amazon Web Services EC2 - Ubuntu</p>
<p>Amazon Web Services RD2 - mySQL</p>
<p>Amazon Web Servies Https LB - load balanced</p>

<h2>View the Entire Project in Action</h2>
<p>It is best to open up both the server as well as the client along with the instructions to see the project in action.</p>
<p>Click <a href="https://kvserver.sozoalvin.com">here</a> for the backend server</p>
<p>Click <a href="https://kvclient.sozoalvin.com">here</a> for the client server</p>

<h2>Understand How RESTAPI works in the Project</h2>
<p>Click<a href="https://github.com/sozoalvin/Food-Delivery-Server/blob/master/Instructions/Understanding_RESTAPI.pdf"> Here </a> to learn more about how RESTAPI was deployed on the backend server as well as the client's</p>

<h2>Instructions on How to Use</h2>
<h4>Click on image for a direct link to the instructions in PDF format.</h4>
<a href="https://github.com/sozoalvin/Food-Delivery-Server/blob/master/Instructions/Instructions%20on%20Navigating%20Backend%20Server%20KV%20Food%20Delivery%20Service.pdf"><img src="https://i.imgur.com/C9XIJuC.png" width="700"></a>

<h2>Development Envirnonment</h2> 
<p>For development purposes in your own environment, use the following command once you've in the working directory</p>

```go run .```

<p>You can access the server on <i>localhost</i> or <i>localhost:80</i> if the former doesn't work</p>

<h2>Testing for Production</h2>
<p>To check if your program is ready for production, use the following command.</p>

```go run . -productionFlg```

<p>You can access the server on <i>localhost</i> if the former doesn't work</p>
<p>When production flag is activated, requests routed on port:80 will be automatically routed to port:443</p>
<p>If you open up chrome's dev tools, you'll see HTTP/1.1 and then see a HTTP/2 almost immedidately</p>

<h2>Launching on an Actual Web Server</h2>

```go run . -productionFlg -domain www.yourdomainname.com```

<p>On an actual web server, your domain name has to be specified. HTTPS will be enabled and cerificate will be provided by; courtesy of let's encrypt</p>

<h2>Search Away to Discover Great Tasting Food!</h2>
<img src = "https://i.imgur.com/XFUVxvV.png" width="700">

<h2>Access Your Own Account's API</h2>
<p>Lost API? Revoke with a simple click and lock out all clients using your old API key</p>
<img src = "https://i.imgur.com/ojn0LQ4.png" width="700">

<h2>Easily Add Any Items to Cart</h2>
<img src = "https://i.imgur.com/FKjFCQ8.png" width="700">

<h2>Admin Users Get Special Settings for Service Recovery</h2>
<img src = "https://i.imgur.com/16XvZTt.png" width="700">

<h2>Regular Account Types at Cart Option</h2>
<img src = "https://i.imgur.com/pHjQuLW.png" width="700">

<h2>Customer's Checkout Page</h2>
<img src = "https://i.imgur.com/PVTs0E7.png" width="700">

<h2>Priority Queues for Service Recvery</h2>
<img src = "https://i.imgur.com/c63vP9z.png" width="700">

<h2>Driver Assignment</h2>
<img src = "https://i.imgur.com/1Nye0iX.png" width="700">
<img src = "https://i.imgur.com/bhKLJWl.png" width="700">

<h2>Read, Get Set, Dispatch and it's a GO!</h2>
<img src = "https://i.imgur.com/xJjm2yX.png" width="700">


<h2>Easily View All System Information</h2>
<img src = "https://i.imgur.com/DCrvoAL.png" width="700">
<img src = "https://i.imgur.com/qhb7rVt.png" width="700">

