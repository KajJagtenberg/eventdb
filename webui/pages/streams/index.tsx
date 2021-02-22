import { getApolloClient } from '../../apollo';
import Layout from '../../components/Layout';
import Navbar from '../../components/Navbar';
import StreamTable from '../../components/StreamTable';
import { title } from '../../vars/title';

const Streams = () => {
  return (
    <Layout title={`${title} - Streams`}>
      <Navbar />

      <StreamTable />
    </Layout>
  );
};

export default Streams;
