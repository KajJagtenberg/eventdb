import { Box } from '@chakra-ui/react';
import Head from 'next/head';

const Layout = ({ title, children }) => {
  return (
    <>
      <Head>
        <title>{title}</title>
      </Head>

      <Box>{children}</Box>
    </>
  );
};

export default Layout;
