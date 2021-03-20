import { Box } from '@chakra-ui/react';
import Head from 'next/head';

const Home = () => {
  return (
    <Box
      style={{
        height: '100vh',
      }}
    >
      <Head>
        <title>EventflowDB</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>
    </Box>
  );
};

export default Home;
