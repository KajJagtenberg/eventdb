import Layout from '../components/Layout';
import Navbar from '../components/Navbar';
import StreamTable from '../components/StreamTable';

const Home = () => {
  return (
    <Layout title="EventDB">
      <Navbar />

      <StreamTable />
    </Layout>
  );
};

export default Home;
