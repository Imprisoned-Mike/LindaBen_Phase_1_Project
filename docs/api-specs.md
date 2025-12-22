# API Specifications

This document outlines the RESTful HTTP API specifications for the web application, detailing endpoints, request/response formats, and error handling.

## Base URL

All API endpoints are prefixed with `/api`. For example, a request to login would be `POST /api/auth/login`.

## Authentication

All protected endpoints require a JWT (JSON Web Token) to be sent in the `Authorization` header as a Bearer token.
Example: `Authorization: Bearer <your_jwt_token>`

## Common Error Responses

| Status Code | Description                                  | Response Body Example                               |
| :---------- | :------------------------------------------- | :-------------------------------------------------- |
| `400 Bad Request`   | Invalid request payload or parameters.       | `{ "message": "Invalid input data", "details": [...] }` |
| `401 Unauthorized`  | Authentication required or invalid token.    | `{ "message": "Authentication required" }`          |
| `403 Forbidden`     | User does not have permission to access.     | `{ "message": "Access denied" }`                    |
| `404 Not Found`     | The requested resource was not found.        | `{ "message": "Resource not found" }`               |
| `422 Unprocessable Entity` | Semantic errors in the request.               | `{ "message": "Validation failed", "errors": [...] }` |
| `500 Internal Server Error` | An unexpected error occurred on the server.  | `{ "message": "Internal server error" }`            |

## Data Models (Types)

### User

Represents a user in the system.

```typescript
export interface User {
  id: number;
  name: string;
  email: string;
  phone: string;
  roles: string; // e.g., "school_admin:123,admin"
  avatarId?: number;
  avatar?: File;
}
```

### LoginResponse

The response returned upon successful user login.

```typescript
export interface LoginResponse {
  token: string;
  user: User;
}
```

### RoleParsed

Represents a parsed user role with an optional entity ID.

```typescript
export interface RoleParsed {
  role: 'admin' | 'school_admin' | 'vendor_admin' | 'user';
  entityId: string | null;
}
```

### PaginationParams

Parameters for paginating results.

```typescript
export interface PaginationParams {
  page?: number;
  pageSize?: number;
}
```

### SortParams

Parameters for sorting results.

```typescript
export interface SortParams {
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}
```

### BaseFilterParams

Base parameters for filtering, combining pagination and sorting.

```typescript
export interface BaseFilterParams extends PaginationParams, SortParams {
  expand?: string[]; // e.g., "expand=school,contact"
}
```

### PaginatedResponse

Standard structure for paginated API responses.

```typescript
export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    total: number;
    totalUnfiltered: number;
    page: number;
    pageSize: number;
    totalPages: number;
  };
}
```

### Order

Represents an order for a delivery.

```typescript
export interface Order {
  id: number;
  item: string;
  quantity: number;
  unitPrice?: number;
  vendorId: number;
  vendor?: Vendor;
  isInternal: boolean;
  completedAt?: string; // ISO date string
  status: 'pending' | 'confirmed' | 'completed' | 'cancelled';
  notes?: string;
}
```

### Delivery

Represents a delivery containing multiple orders to a school.

```typescript
export interface Delivery {
  id: number;
  orders: Order[];
  contract: 'hold' | 'completed';
  schoolId: number;
  school?: School; // Expanded school object
  packageType: 'standard' | 'premium' | 'standard-holiday' | 'premium-holiday';
  notes?: string;
  scheduledAt?: string; // ISO date string
}
```

### DeliveryChangeLog

Log entry for changes made to a delivery.

```typescript
export interface DeliveryChangeLog {
  id: number;
  deliveryId: number;
  changedByUserId: number;
  changedByUser?: User; // Expanded user object
  changedAt: string; // ISO date string
  fieldName: 'scheduledAt' | 'packageType' | 'notes' | 'contract' | 'schoolId';
  oldValue: string | number | boolean | null | undefined;
  newValue: string | number | boolean | null | undefined;
}
```

### OrderChangeLog

Log entry for changes made to an order.

```typescript
export interface OrderChangeLog {
  id: number;
  orderId: number;
  changedByUserId: number;
  changedByUser?: User; // Expanded user object
  changedAt: string; // ISO date string
  fieldName: 'status' | 'quantity' | 'item' | 'unitPrice' | 'notes' | 'isInternal' | 'vendorId';
  oldValue: string | number | boolean | null | undefined;
  newValue: string | number | boolean | null | undefined;
}
```

### File

Represents a file uploaded to the system.

```typescript
export interface File {
  id: number;
  url: string;
}
```

### Coordinate

Geographic coordinates (latitude and longitude).

```typescript
export interface Coordinate {
  lat: number;
  lng: number;
}
```

### School

Represents a school entity.

```typescript
export interface School {
  id: number;
  name: string;
  address: string;
  coordinate: Coordinate;
  contactId?: number;
  contact?: User; // Expanded user object
}
```

### Vendor

Represents a vendor entity.

```typescript
export type VendorType = 'produce' | 'shelf_stable' | 'packaging';

export interface Vendor {
  id: number;
  name: string;
  address: string;
  coordinate: Coordinate;
  type: VendorType;
}
```

## API Endpoints

### Auth API

#### `POST /api/auth/login`

Authenticates a user and returns a JWT.

*   **Request Body:**
    ```json
    {
      "email": "string",
      "password": "string"
    }
    ```
*   **Success Response:** `200 OK`
    *   Body: `LoginResponse`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`

#### `POST /api/auth/logout`

Invalidates the user's session/token.

*   **Success Response:** `204 No Content`
*   **Error Responses:** `401 Unauthorized`

#### `POST /api/auth/impersonate`

Allows an admin to impersonate another user.

*   **Request Body:**
    ```json
    {
      "userId": 123
    }
    ```
*   **Success Response:** `200 OK`
    *   Body: `LoginResponse`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

### Users API

#### `GET /api/users`

Retrieves a paginated list of users.

*   **Query Parameters:**
    *   `search` (string, optional): Free-text search across name, email, and phone.
    *   `role` (`RoleParsed["role"]`, optional): Filter by a specific parsed role.
    *   `entityId` (string, optional): Filter by the entity ID embedded in roles.
    *   `id` (number, optional): Filter by a specific user ID.
    *   `hasRole` (string, optional): Filter by users possessing a specific role.
    *   `email` (string, optional): Filter by email.
    *   `name` (string, optional): Filter by name.
    *   `page` (number, optional): Page number for pagination (default: 1).
    *   `pageSize` (number, optional): Number of items per page (default: 10).
    *   `sortBy` (string, optional): Field to sort by.
    *   `sortOrder` ('asc' | 'desc', optional): Sort order (default: 'asc').
    *   `expand` (string, optional): Comma-separated list of fields to expand (e.g., `avatar`).
*   **Success Response:** `200 OK`
    *   Body: `PaginatedResponse<User>`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`

#### `GET /api/users/{id}`

Retrieves a single user by ID.

*   **Path Parameters:**
    *   `id` (number): The ID of the user.
*   **Query Parameters:**
    *   `expand` (string, optional): Comma-separated list of fields to expand (e.g., `avatar`).
*   **Success Response:** `200 OK`
    *   Body: `User`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `POST /api/users`

Creates a new user.

*   **Request Body:** `Omit<User, 'id' | 'avatarId' | 'avatar'>` (excluding generated fields)
*   **Success Response:** `201 Created`
    *   Body: `User`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `422 Unprocessable Entity`

#### `PUT /api/users/{id}`

Updates an existing user.

*   **Path Parameters:**
    *   `id` (number): The ID of the user to update.
*   **Request Body:** `Partial<User>`
*   **Success Response:** `200 OK`
    *   Body: `User`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `422 Unprocessable Entity`

#### `DELETE /api/users/{id}`

Deletes a user.

*   **Path Parameters:**
    *   `id` (number): The ID of the user to delete.
*   **Success Response:** `204 No Content`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `POST /api/users/{id}/avatar`

Uploads an avatar for a user.

*   **Path Parameters:**
    *   `id` (number): The ID of the user.
*   **Request Body:** `multipart/form-data` with a `file` field containing the image file.
*   **Success Response:** `200 OK`
    *   Body: `{ "url": "string" }` (URL of the uploaded avatar)
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

### Schools API

#### `GET /api/schools`

Retrieves a paginated list of schools.

*   **Query Parameters:**
    *   `search` (string, optional): Free-text search across name and address.
    *   `hasContact` (boolean, optional): Filter by schools that have a contact user assigned.
    *   `contactUserId` (number, optional): Filter by a specific contact user ID.
    *   `page` (number, optional): Page number for pagination (default: 1).
    *   `pageSize` (number, optional): Number of items per page (default: 10).
    *   `sortBy` (string, optional): Field to sort by.
    *   `sortOrder` ('asc' | 'desc', optional): Sort order (default: 'asc').
    *   `expand` (string, optional): Comma-separated list of fields to expand (e.g., `contact`).
*   **Success Response:** `200 OK`
    *   Body: `PaginatedResponse<School>`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`

#### `GET /api/schools/{id}`

Retrieves a single school by ID.

*   **Path Parameters:**
    *   `id` (number): The ID of the school.
*   **Query Parameters:**
    *   `expand` (string, optional): Comma-separated list of fields to expand (e.g., `contact`).
*   **Success Response:** `200 OK`
    *   Body: `School`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `POST /api/schools`

Creates a new school.

*   **Request Body:** `Omit<School, 'id' | 'contact'>`
*   **Success Response:** `201 Created`
    *   Body: `School`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `422 Unprocessable Entity`

#### `PUT /api/schools/{id}`

Updates an existing school.

*   **Path Parameters:**
    *   `id` (number): The ID of the school to update.
*   **Request Body:** `Partial<School>`
*   **Success Response:** `200 OK`
    *   Body: `School`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `422 Unprocessable Entity`

#### `DELETE /api/schools/{id}`

Deletes a school.

*   **Path Parameters:**
    *   `id` (number): The ID of the school to delete.
*   **Success Response:** `204 No Content`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

### Vendors API

#### `GET /api/vendors`

Retrieves a paginated list of vendors.

*   **Query Parameters:**
    *   `search` (string, optional): Free-text search across vendor name and address.
    *   `types` (`VendorType[]`, optional): Filter by one or more vendor categories (e.g., `types=produce,packaging`).
    *   `page` (number, optional): Page number for pagination (default: 1).
    *   `pageSize` (number, optional): Number of items per page (default: 10).
    *   `sortBy` (string, optional): Field to sort by.
    *   `sortOrder` ('asc' | 'desc', optional): Sort order (default: 'asc').
*   **Success Response:** `200 OK`
    *   Body: `PaginatedResponse<Vendor>`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`

#### `GET /api/vendors/{id}`

Retrieves a single vendor by ID.

*   **Path Parameters:**
    *   `id` (number): The ID of the vendor.
*   **Success Response:** `200 OK`
    *   Body: `Vendor`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `POST /api/vendors`

Creates a new vendor.

*   **Request Body:** `Omit<Vendor, 'id'>`
*   **Success Response:** `201 Created`
    *   Body: `Vendor`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `422 Unprocessable Entity`

#### `PUT /api/vendors/{id}`

Updates an existing vendor.

*   **Path Parameters:**
    *   `id` (number): The ID of the vendor to update.
*   **Request Body:** `Partial<Vendor>`
*   **Success Response:** `200 OK`
    *   Body: `Vendor`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `422 Unprocessable Entity`

#### `DELETE /api/vendors/{id}`

Deletes a vendor.

*   **Path Parameters:**
    *   `id` (number): The ID of the vendor to delete.
*   **Success Response:** `204 No Content`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

### Deliveries API

#### `GET /api/deliveries`

Retrieves a paginated list of deliveries.

*   **Query Parameters:**
    *   `search` (string, optional): Free-text search across school name or notes.
    *   `schoolId` (number, optional): Filter by a specific school ID.
    *   `scheduledFrom` (string, optional, ISO date string): Filter deliveries scheduled from this date.
    *   `scheduledTo` (string, optional, ISO date string): Filter deliveries scheduled to this date.
    *   `contract` (`Delivery["contract"]` or `Delivery["contract"][]`, optional): Filter by contract type(s).
    *   `packageType` (`Delivery["packageType"]` or `Delivery["packageType"][]`, optional): Filter by package type(s).
    *   `status` (`Order["status"]` or `Order["status"][]`, optional): Filter by overall computed order status within a delivery.
    *   `page` (number, optional): Page number for pagination (default: 1).
    *   `pageSize` (number, optional): Number of items per page (default: 10).
    *   `sortBy` (string, optional): Field to sort by.
    *   `sortOrder` ('asc' | 'desc', optional): Sort order (default: 'asc').
    *   `expand` (string, optional): Comma-separated list of fields to expand (e.g., `school,orders.vendor`).
*   **Success Response:** `200 OK`
    *   Body: `PaginatedResponse<Delivery>`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`

#### `GET /api/deliveries/{id}`

Retrieves a single delivery by ID.

*   **Path Parameters:**
    *   `id` (number): The ID of the delivery.
*   **Query Parameters:**
    *   `expand` (string, optional): Comma-separated list of fields to expand (e.g., `school,orders.vendor`).
*   **Success Response:** `200 OK`
    *   Body: `Delivery`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `POST /api/deliveries`

Creates a new delivery.

*   **Request Body:** `Omit<Delivery, 'id' | 'school'>`
*   **Success Response:** `201 Created`
    *   Body: `Delivery`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `422 Unprocessable Entity`

#### `PUT /api/deliveries/{id}`

Updates an existing delivery.

*   **Path Parameters:**
    *   `id` (number): The ID of the delivery to update.
*   **Request Body:** `Partial<Delivery>`
*   **Success Response:** `200 OK`
    *   Body: `Delivery`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `422 Unprocessable Entity`

#### `DELETE /api/deliveries/{id}`

Deletes a delivery.

*   **Path Parameters:**
    *   `id` (number): The ID of the delivery to delete.
*   **Success Response:** `204 No Content`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `GET /api/deliveries/{id}/logs`

Retrieves the change logs for a specific delivery.

*   **Path Parameters:**
    *   `id` (number): The ID of the delivery.
*   **Success Response:** `200 OK`
    *   Body: `DeliveryChangeLog[]`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `POST /api/deliveries/{deliveryId}/orders`

Adds a new order to a specific delivery.

*   **Path Parameters:**
    *   `deliveryId` (number): The ID of the delivery.
*   **Request Body:** `Omit<Order, 'id' | 'vendor'>`
*   **Success Response:** `201 Created`
    *   Body: `Delivery` (the updated delivery object)
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `422 Unprocessable Entity`

#### `DELETE /api/deliveries/{deliveryId}/orders/{orderId}`

Removes an order from a specific delivery.

*   **Path Parameters:**
    *   `deliveryId` (number): The ID of the delivery.
    *   `orderId` (number): The ID of the order to remove.
*   **Success Response:** `204 No Content`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

### Orders API

#### `GET /api/orders/{id}`

Retrieves a single order by ID, including its scheduled delivery time.

*   **Path Parameters:**
    *   `id` (number): The ID of the order.
*   **Query Parameters:**
    *   `expand` (string, optional): Comma-separated list of fields to expand (e.g., `vendor`).
*   **Success Response:** `200 OK`
    *   Body: `OrderWithSchedule` (`Order` with an optional `scheduledAt` field)
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `PUT /api/orders/{id}`

Updates an existing order.

*   **Path Parameters:**
    *   `id` (number): The ID of the order to update.
*   **Request Body:** `Partial<Order>`
*   **Success Response:** `200 OK`
    *   Body: `Order`
*   **Error Responses:** `400 Bad Request`, `401 Unauthorized`, `403 Forbidden`, `404 Not Found`, `422 Unprocessable Entity`

#### `GET /api/orders/{id}/logs`

Retrieves the change logs for a specific order.

*   **Path Parameters:**
    *   `id` (number): The ID of the order.
*   **Success Response:** `200 OK`
    *   Body: `OrderChangeLog[]`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `DELETE /api/orders/{id}`

Deletes an order.

*   **Path Parameters:**
    *   `id` (number): The ID of the order to delete.
*   **Success Response:** `204 No Content`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

#### `GET /api/orders/items/search`

Searches for available items.

*   **Query Parameters:**
    *   `query` (string): The search query for item names.
*   **Success Response:** `200 OK`
    *   Body: `Array<{ name: string; defaultPrice: number }>`
*   **Error Responses:** `401 Unauthorized`, `403 Forbidden`