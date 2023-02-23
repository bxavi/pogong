// Code generated by tygo. DO NOT EDIT.
"// Generated by Tygo"
"// https://golangexample.com/generate-typescript-types-from-golang-source-code/"

//////////
// source: models.go

export interface Account {
  email: string;
  password: string;
  created_at: string /* RFC3339 */;
}

//////////
// source: query.sql.go

export interface CreateAccountParams {
  email: string;
  password: string;
}
export interface ListAccountParams {
  offset: any /* sql.NullInt32 */;
  limit: any /* sql.NullInt32 */;
}
export interface UpdateAccountParams {
  id: number /* int64 */;
  email: string;
  password: string;
}

//////////
// source: store.go

/**
 * Store defines all functions to execute db queries and transactions
 */
export type Store = 
    Querier;
/**
 * SQLStore provides all functions to execute SQL queries and transactions
 */
export interface SQLStore {
}
