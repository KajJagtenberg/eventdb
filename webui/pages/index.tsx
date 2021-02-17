import { useQuery } from 'react-query';

import HomeOverview from '../components/HomeOverview';
import Layout from '../components/Layout';
import Navbar from '../components/Navbar';
import { backend } from '../vars/backend';

const fetchInfo = async () => {
  const response = await fetch(`${backend}`);
  return response.json();
};

const Home = () => {
  const { data } = useQuery('info', fetchInfo);
  const { version, size, human_size, event_count, stream_count } = data || {};

  return (
    <Layout title="EventDB - Home">
      <Navbar />

      <HomeOverview
        version={version}
        size={size}
        human_size={human_size}
        event_count={event_count}
        stream_count={stream_count}
      />
    </Layout>
  );
};

export default Home;
