// Code generated by tygo. DO NOT EDIT.
"// Generated by Tygo"
"// https://golangexample.com/generate-typescript-types-from-golang-source-code/"

//////////
// source: models.go

export interface Account {
  email: string;
  password: string;
}
export interface Test {
  field: any /* nulls.String */;
}

//////////
// source: query.sql.go

export interface CreateAccountsParams {
  email: string;
  password: string;
}
export interface CreateTestParams {
  id: any /* nulls.Int */;
  field: any /* nulls.String */;
}
export interface UpdateAccountsParams {
  id: number /* int64 */;
  email: string;
  password: string;
}
