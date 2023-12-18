# License Server and Manager

The goal of this project was to create a software licensing server with a REST API for users, along with a CLI for admin management.

### License Generation

The CLI can be used to manage licenses; additionally, it provides functionality to generate empty licenses whose product keys can be distributed to consumers. Before a license can be generated the server requires two resources to exist, a DB entity corresponding to the product and version the license will be attached to and an ed25519 key pair associated with said product/version. Once these have been created a new license can be generated per the workflow below:

1. License metadata is initialized as a JSON string
2. Metadata JSON is hashed and signed with the product/version private key
3. Signature is AES encrypted and base58 encoded for distribution as the product key
4. License is stored in db as `{"key": encrypted signature, "attributes": license metadata}`

### License Issuing

Using the REST API, clients corresponding to products in the server database can POST consumer product keys to validate the keys and issue licenses. When it receives a request to the `/issue` endpoint, the server completes the following workflow:

1. DB key lookup and signature verification if key is found & unassigned
2. Key is bound to the client machine or all machines if none provided
3. Client activation response is sent: `{"identity": "", "machine": "", "issue_date": "", "expiration_date": "", "refresh_date": ""}`

On receipt of the activation data, the client should store it locally and repeat the key verification process before the `refresh_date` to avoid invalidating the key.

### API Docs

[https://julianstephens.github.io/license-server/v3/](https://julianstephens.github.io/license-server/v3/)

### Future Work

- Management GUI
- Client libraries
