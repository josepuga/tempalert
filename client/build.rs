fn main() {
    prost_build::compile_protos(&["proto/temperatures.proto"], &["proto/"]).unwrap();
}

