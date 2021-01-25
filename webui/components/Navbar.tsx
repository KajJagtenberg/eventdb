import Link from 'next/link';
import { useRouter } from 'next/router';

import { Box, Flex, Heading, Link as UILink, Stack } from '@chakra-ui/react';

const NavbarLink = ({ href, title }) => {
  const router = useRouter();

  return (
    <Box
      bg={router.pathname === href ? 'teal.600' : 'transparent'}
      py={1}
      px={2}
      rounded="md"
    >
      <Link href={href}>
        <UILink fontWeight="semibold" color="white">
          {title}
        </UILink>
      </Link>
    </Box>
  );
};

const Navbar = () => {
  return (
    <Flex
      p={4}
      boxShadow="md"
      bg="teal.500"
      direction="row"
      justifyContent="space-between"
      alignItems="center"
    >
      <Heading size="md" color="white">
        EventDB
      </Heading>

      <Stack direction="row">
        <NavbarLink title="Streams" href="/" />
        <NavbarLink title="Projections" href="/projections" />
        <NavbarLink title="Users" href="/users" />
      </Stack>
    </Flex>
  );
};

export default Navbar;
