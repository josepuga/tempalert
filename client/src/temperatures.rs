//use prost::Message;
// @generated
// This file is @generated by prost-build.
#[derive(Clone, PartialEq, ::prost::Message)]
//#[derive(Clone, PartialEq, Message)]
pub struct SensorAlert {
    #[prost(int32, tag="1")]
    pub sensor_id: i32,
    #[prost(int32, tag="2")]
    pub temperature: i32,
    #[prost(string, tag="3")]
    pub timestamp: ::prost::alloc::string::String,
}
// @@protoc_insertion_point(module)
