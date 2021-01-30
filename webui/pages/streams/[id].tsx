import { useRouter } from 'next/router';

import EventsTable from '../../components/EventsTable';
import Layout from '../../components/Layout';
import Navbar from '../../components/Navbar';

const Stream = () => {
  const router = useRouter();
  const { id: stream } = router.query;

  return (
    <Layout title={`EventDB - ${stream}`}>
      <Navbar />

      <EventsTable stream={stream} />
    </Layout>
  );
};

export default Stream;
