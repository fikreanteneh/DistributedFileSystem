import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect, useState } from 'react';
import './App.css';



function App() {
  const [file, setFile] = useState(null);
  const [myFiles, setMyFiles] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [downloaded, setDownloaded] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      console.log("use effect==========")
      try {
        const response = await fetch('http://localhost:8000/view', 
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        const data = await response.json();
        console.log("data", data)
        if (data.files) {
          
          setMyFiles(data.files);
        }
      } catch (error) {
        console.error('Error:', error);
      }
    }
    fetchData();
  }, []);



  const handleFileChange = (event) => {
    const selectedFile = event.target.files[0];
    setFile(selectedFile);
    console.log("selectedFile" , selectedFile)
  };

  const handleUpload = () => {
    setUploading(true);
    if (file) {
      initiateUpload();
      setUploading(false);
    } else {
      alert('Please select a file first.');
    }
  };


  const handleDownload = async (filee) => {
    try {
      setDownloaded(true);
      const response = await fetch(`http://localhost:8000/get?id=${filee.fileIdentifier}`)
      const data = await response.json();
      const locations = data.Locations;
      const fileName = data.FileName;

      console.log("locations", locations)

      const file =  await getChunks(locations, filee.fileIdentifier, filee.fileName);
      console.log("file generated", filee)
      setDownloaded(false);

    }catch (e){
      console.log(e, "error")
    }

  }
  
  const initiateUpload = async () => {
    console.log("file" , file)
    try {
      // Step 1: Request metadata from the master
      const response = await fetch('http://localhost:8000/upload', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          fileName: file.name,
          fileSize: file.size, // Set the actual file size if available
        }),
      });
    

      if (!response.ok) {
        throw new Error(`Failed to initiate upload: ${response.statusText}`);
      }

      const uploadInitResponse = await response.json();

      // Step 2: Upload the file to chunk servers
      await uploadToChunkServers(uploadInitResponse);

      // Optionally, you can perform additional actions after successful upload
      console.log('File uploaded successfully!');
    } catch (error) {
      // this.setState({ error: error.message || 'An error occurred.' });
      console.log("error" , error)
    }
  };



  const uploadToChunkServers = async (uploadInitResponse) => {
    try {

      for (let i = 0; i < uploadInitResponse.numberOfChunks; i++) {
        const partSize = Math.min(uploadInitResponse.chunkSize, file.size - i * uploadInitResponse.chunkSize);
        const partBuffer = new Uint8Array(partSize);

        // Read the chunk from the file
        const chunk = await readFileChunk(file, i * uploadInitResponse.chunkSize, partSize);
        partBuffer.set(new Uint8Array(chunk));

        for (let k = 0; k < uploadInitResponse.chunkservers.length; k++) {
          const randomChunkserver = uploadInitResponse.chunkservers[k];
          await uploadChunk(randomChunkserver, partBuffer, `${uploadInitResponse.identifier}_${i}`);
        }

      }
    } catch (error) {
      throw new Error(`Error uploading to chunk servers: ${error.message}`);
    }
  };

  const readFileChunk = (file, start, length) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onloadend = () => {
        resolve(reader.result);
      };

      reader.onerror = () => {
        reject(new Error('Error reading file chunk.'));
      };

      const blob = file.slice(start, start + length);
      reader.readAsArrayBuffer(blob);
    });
  };

  const uploadChunk = async (chunkserverUrl, chunk, chunkIdentifier) => {
    try {
      const formData = new FormData();
      formData.append('chunk', new Blob([chunk]), chunkIdentifier);

      const response = await fetch(`http://${chunkserverUrl}uploadChunk`, {
        method: 'POST',
        body: formData,
      });

      if (response.status !== 200) {
        throw new Error(`Failed to upload chunk to ${chunkserverUrl}: ${response.statusText}`);
      }

      // Optionally, you can perform additional actions after successful upload
      console.log(`Chunk uploaded successfully to ${chunkserverUrl}`);
    } catch (error) {
      throw new Error(`Error uploading chunk: ${error.message}`);
    }
  };


  const initGet = async () => {
    const url = `${client.Master}get?id=${getIdentifierFromFilename(client.Filename)}`;
    
    try {
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`Failed to initiate get: ${response.statusText}`);
        }

        const getResponse = await response.json();
        return getResponse;
    } catch (error) {
        console.error("Error in initGet:", error);
        throw error;
    }
}


// function getIdentifierFromFilename(filename) {
//     // const hashed = crypto.createHash('sha256').update(filename).digest('hex');
//     const hashed = crypto.SHA256(filename).toString();
//     console.log("hashed", hashed)
//     return hashed;
// }


const getChunks = async (locations, id, Filename) => {
  const chunks = [];

  try {
      for (let i = 0; i < locations.length; i++) {
          const chunkIdentifier = `${id}_${i}`;
          const chunkserver = locations[i];

          const chunk = await getChunk(chunkIdentifier, chunkserver);
          console.log("chunk", chunk);
          chunks.push(chunk);
      }
      console.log("chunks", chunks);

      // Concatenate the chunks into a single Uint8Array
      const concatenatedChunks = new Uint8Array(chunks.reduce((acc, chunk) => acc.concat(Array.from(chunk)), []));

      // Create a Blob from the concatenated chunks
      const blob = new Blob([concatenatedChunks], { type: 'application/octet-stream' });

      // Create a download link and trigger a click event to download the file
      const downloadLink = document.createElement('a');
      downloadLink.href = URL.createObjectURL(blob);
      downloadLink.download = Filename;
      downloadLink.click();
  } catch (error) {
      console.error("Error in getChunks:", error);
      throw error;
  }
};


const getChunk = async (chunkIdentifier, chunkserver) => {
    const url = `http://${chunkserver}get?id=${chunkIdentifier}`;

    try {
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`Failed to get chunk: ${response.statusText}`);
        }

        const chunk = await response.arrayBuffer();
        return new Uint8Array(chunk);
    } catch (error) {
        console.error("Error in getChunk:", error);
        throw error;
    }
}


  return (
    <div className="container mt-5">
      <div className="upload-container">
        <h1 className="mb-4">File Upload App</h1>
        <label htmlFor="file-input" className="file-label">
          <div className={`file-picker ${file ? 'file-selected' : ''}`}>
            <span>{file ? 'File Selected' : 'Upload File'}</span>
            <input id="file-input" type="file" onChange={handleFileChange} className="form-control-file" />
          </div>
        </label>
      </div>
        {file && (
          <button onClick={handleUpload} className="btn btn-dark mt-3">
            Upload File
          </button>
        )}

      {/* Display the list of uploaded files in a table-like structure */}
      <div className="mt-5">
      <h2>My Files</h2>
      <table className="table table-bordered" style={{ maxWidth: '800px' }}>
        <thead className="thead-dark">
          <tr>
            <th scope="col" style={{ width: '10%' }}>Row No</th>
            <th scope="col" style={{ width: '75%' }}>Filename</th>
            <th scope="col" style={{ width: '15%', textAlign: 'right' }}>Action</th>
          </tr>
        </thead>
        <tbody>
          {myFiles.map((file, index) => (
            <tr key={file.id}>
              <td>{index + 1}</td>
              <td className="text-left">{file.fileName}</td>
              <td className="text-right">
                <button onClick={() => handleDownload(file)} className="btn btn-success">
                  Download
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      </div>
      {uploading && <h3>Uh3loading...</h3>}
      {downloaded && <h3>Downloading...</h3>}

    </div>
  );
};

export default App;