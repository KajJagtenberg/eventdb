import Layout from '../../components/Layout';
import Navbar from '../../components/Navbar';
import StreamTable from '../../components/StreamTable';

const Streams = () => {
  return (
    <Layout title="EventDB - Streams">
      <Navbar />

      <StreamTable />
    </Layout>
  );
};

export default Streams;
