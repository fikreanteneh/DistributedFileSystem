// /* eslint-disable react/prop-types */
// import React, { } from 'react';

// class FileClient extends React.Component {
//   constructor(props) {
//     super(props);

//     this.state = {
//       master: props.master,
//       action: props.action,
//       filename: props.filename,
//       fileSize: props.fileSize,
//       outputFilename: props.outputFilename,
//       uploadInitResponse: null,
//       error: null,
//     };
//   }

//   componentDidMount() {
//     this.initiateUpload();
//   }

//   initiateUpload = async () => {
//     try {
//       // Step 1: Request metadata from the master
//       const response = await fetch(`${master}upload`, {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//         },
//         body: JSON.stringify({
//             fileName: file.name,
//             fileSize: file.size
//         }),
//       });


//       console.log("response" , response)
//       if (!response.ok) {
//         throw new Error(`Failed to initiate upload: ${response.statusText}`);
//       }

//       const uploadInitResponse = await response.json();
//       this.setState({ uploadInitResponse });

//       // Step 2: Upload the file to chunk servers
//       await this.uploadToChunkServers(uploadInitResponse);

//       // Optionally, you can perform additional actions after successful upload
//       console.log('File uploaded successfully!');
//     } catch (error) {
//       this.setState({ error: error.message || 'An error occurred.' });
//     }
//   };

//   uploadToChunkServers = async (uploadInitResponse) => {
//     try {
//       // const { filename } = this.state;

//       // Assuming you have the file input in your React component
//       const fileInput = document.getElementById('fileInput');
//       const file = fileInput.files[0];

//       if (!file) {
//         throw new Error('No file selected.');
//       }

//       for (let i = 0; i < uploadInitResponse.numberOfChunks; i++) {
//         const partSize = Math.min(uploadInitResponse.chunkSize, file.size - i * uploadInitResponse.chunkSize);
//         const partBuffer = new Uint8Array(partSize);

//         // Read the chunk from the file
//         const chunk = await this.readFileChunk(file, i * uploadInitResponse.chunkSize, partSize);
//         partBuffer.set(new Uint8Array(chunk));

//         // Upload the chunk to a random chunk server
//         const randomChunkserver = uploadInitResponse.chunkservers[Math.floor(Math.random() * uploadInitResponse.chunkservers.length)];
//         await this.uploadChunk(randomChunkserver, partBuffer, `${uploadInitResponse.identifier}_${i}`);
//       }
//     } catch (error) {
//       throw new Error(`Error uploading to chunk servers: ${error.message}`);
//     }
//   };

//   readFileChunk = (file, start, length) => {
//     return new Promise((resolve, reject) => {
//       const reader = new FileReader();

//       reader.onloadend = () => {
//         resolve(reader.result);
//       };

//       reader.onerror = () => {
//         reject(new Error('Error reading file chunk.'));
//       };

//       const blob = file.slice(start, start + length);
//       reader.readAsArrayBuffer(blob);
//     });
//   };

//   uploadChunk = async (chunkserverUrl, chunk, chunkIdentifier) => {
//     try {
//       const formData = new FormData();
//       formData.append('chunk', new Blob([chunk]), chunkIdentifier);

//       const response = await fetch(`http://${chunkserverUrl}uploadChunk`, {
//         method: 'POST',
//         body: formData,
//       });

//       if (response.status !== 200) {
//         throw new Error(`Failed to upload chunk to ${chunkserverUrl}: ${response.statusText}`);
//       }

//       // Optionally, you can perform additional actions after successful upload
//       console.log(`Chunk uploaded successfully to ${chunkserverUrl}`);
//     } catch (error) {
//       throw new Error(`Error uploading chunk: ${error.message}`);
//     }
//   };

//   render() {
//     const { error } = this.state;

//     return (
//       <div>
//         {error && <p>Error: {error}</p>}
//         {/* Add any other UI elements or indicators as needed */}
//       </div>
//     );
//   }
// }

// export default FileClient;
