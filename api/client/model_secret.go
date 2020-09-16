/*
 * MLP API
 *
 * API Guide for accessing MLP API
 *
 * API version: 0.4.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

type Secret struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}
