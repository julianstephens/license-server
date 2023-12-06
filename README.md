# license-server

License server and manager that grants access tokens to users/groups.

### DB Schema

#### users

- id: string
- username: string

#### groups

- id: string
- users: string[]

#### products

- id: string
- name: string

#### licenses

- id: string
- productID: string
- uid?: string
- groupID?: string
